package main

import (
	// stdlib
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/run"
	uuid "github.com/satori/go.uuid"
	"go.opencensus.io/plugin/ochttp"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/database/sqlite"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/transport/pb"
	svcevent "github.com/basvanbeek/opencensus-gokit-example/services/event/transport/twirp"
	"github.com/basvanbeek/opencensus-gokit-example/shared/network"
	"github.com/basvanbeek/opencensus-gokit-example/shared/opencensus"
)

func main() {
	var (
		err      error
		instance = uuid.Must(uuid.NewV4())
	)

	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", event.ServiceName,
			"instance", instance,
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	// initialize our OpenCensus configuration
	defer opencensus.Setup(event.ServiceName).Close()

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create our etcd client for Service Discovery
	//
	// we could have used the v3 client but then we must vendor or suffer the
	// following issue originating from gRPC init:
	// panic: http: multiple registrations for /debug/requests
	var sdc etcd.Client
	{
		// create our Go kit etcd client
		sdc, err = etcd.NewClient(ctx, []string{"http://localhost:2379"}, etcd.ClientOptions{})
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	// Create our DB Connection Driver
	var db *sqlx.DB
	{
		db, err = sqlx.Open("sqlite3", "event.db")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		// make sure the DB is in WAL mode
		if _, err = db.Exec(`PRAGMA journal_mode=wal`); err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	// Create our Event Service
	var svc event.Service
	{
		repository, err := sqlite.New(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = implementation.NewService(repository, logger)
		// add service level middlewares here
	}

	// Create our Go kit endpoints for the Event Service
	// var endpoints transport.Endpoints
	// {
	// 	endpoints = transport.MakeEndpoints(svc)
	// 	// add endpoint level middlewares here
	// }

	// run.Group manages our goroutine lifecycles
	// see: https://www.youtube.com/watch?v=LHe1Cb_Ud_M&t=15m45s
	var g run.Group
	{
		// set-up our ZPages handler
		opencensus.ZPages(g, logger)
	}
	{
		// set-up our twirp transport
		var (
			bindIP, _    = network.HostIP()
			eventService = svcevent.NewTwirpServer(svc, logger)
			listener, _  = net.Listen("tcp", bindIP+":0") // dynamic port assignment
			svcInstance  = "/services/" + event.ServiceName + "/twirp/" + instance.String()
			addr         = "http://" + listener.Addr().String()
			ttl          = etcd.NewTTLOption(3*time.Second, 10*time.Second)
			service      = etcd.Service{Key: svcInstance, Value: addr, TTL: ttl}
			registrar    = etcd.NewRegistrar(sdc, service, logger)
			twirpHandler = pb.NewEventServer(eventService, nil)
			router       = mux.NewRouter()
		)

		router.PathPrefix(pb.EventPathPrefix).Handler(twirpHandler)

		// add default ochttp handler for TWIRP
		handler := &ochttp.Handler{
			Handler: router,
		}

		g.Add(func() error {
			registrar.Register()
			return http.Serve(listener, handler)
		}, func(error) {
			registrar.Deregister()
			listener.Close()
		})
	}
	{
		// set-up our signal handler
		var (
			cancelInterrupt = make(chan struct{})
			c               = make(chan os.Signal, 2)
		)
		defer close(c)

		g.Add(func() error {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	// spawn our goroutines and wait for shutdown
	level.Error(logger).Log("exit", g.Run())
}
