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
	"github.com/oklog/run"
	uuid "github.com/satori/go.uuid"
	"go.opencensus.io/plugin/ochttp"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
	devclient "github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/clients/device"
	evtclient "github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/clients/event"
	qrclient "github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/clients/qr/grpc"
	svchttp "github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
	"github.com/basvanbeek/opencensus-gokit-example/shared/network"
	"github.com/basvanbeek/opencensus-gokit-example/shared/oc"
)

func main() {
	var (
		err      error
		instance = uuid.Must(uuid.NewV4())
	)

	// initialize our OpenCensus configuration and defer a clean-up
	defer oc.Setup(frontend.ServiceName).Close()

	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", frontend.ServiceName,
			"instance", instance,
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

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

	var svc frontend.Service
	{
		// create an instancer for the event client
		evtInstancer, err := etcd.NewInstancer(sdc, "/services/"+event.ServiceName+"/twirp", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}
		httpClient := &http.Client{Transport: &ochttp.Transport{}}
		evtClient := evtclient.NewTwirp(evtInstancer, httpClient, logger)

		// create an instancer for the device client
		devInstancer, err := etcd.NewInstancer(sdc, "/services/"+device.ServiceName+"/http", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}

		// initialize our Device client using http transport
		devClient := devclient.NewHTTP(devInstancer, logger)

		// create an instancer for the QR client
		qrInstancer, err := etcd.NewInstancer(sdc, "/services/"+qr.ServiceName+"/grpc", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}
		// initialize QR client
		qrClient := qrclient.New(qrInstancer, logger)

		// create our frontend service
		svc = implementation.NewService(evtClient, devClient, qrClient, logger)
		// add service level middlewares here
	}

	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(svc)

		// trace our server side endpoints
		endpoints = transport.Endpoints{
			Login:        oc.ServerEndpoint("Login")(endpoints.Login),
			EventCreate:  oc.ServerEndpoint("EventCreate")(endpoints.EventCreate),
			EventGet:     oc.ServerEndpoint("EventGet")(endpoints.EventGet),
			EventUpdate:  oc.ServerEndpoint("EventUpdate")(endpoints.EventUpdate),
			EventDelete:  oc.ServerEndpoint("EventDelete")(endpoints.EventDelete),
			EventList:    oc.ServerEndpoint("EventList")(endpoints.EventList),
			UnlockDevice: oc.ServerEndpoint("UnlockDevice")(endpoints.UnlockDevice),
			GenerateQR:   oc.ServerEndpoint("GenerateQR")(endpoints.GenerateQR),
		}
	}

	// run.Group manages our goroutine lifecycles
	// see: https://www.youtube.com/watch?v=LHe1Cb_Ud_M&t=15m45s
	var g run.Group
	{
		// set-up our ZPages handler
		oc.ZPages(g, logger)
	}
	{
		// set-up our http transport
		var (
			bindIP, _   = network.HostIP()
			listener, _ = net.Listen("tcp", bindIP+":0") // dynamic port assignment
			svcInstance = fmt.Sprintf("/services/%s/http/%s/", frontend.ServiceName, instance)
			addr        = "http://" + listener.Addr().String()
			ttl         = etcd.NewTTLOption(3*time.Second, 10*time.Second)
			service     = etcd.Service{Key: svcInstance, Value: addr, TTL: ttl}
			registrar   = etcd.NewRegistrar(sdc, service, logger)
			feService   = svchttp.NewHTTPHandler(endpoints, logger)
		)

		g.Add(func() error {
			registrar.Register()
			return http.Serve(listener, feService)
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
