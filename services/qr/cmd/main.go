package main

import (
	// stdlib
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcd"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/kevinburke/go.uuid"
	"github.com/oklog/run"
	"google.golang.org/grpc"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport"
	grpctransport "github.com/basvanbeek/opencensus-gokit-example/services/qr/transport/grpc"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/network"
	"github.com/basvanbeek/opencensus-gokit-example/shared/oc"
)

func main() {
	var (
		err      error
		instance = uuid.NewV4()
	)

	// initialize our OpenCensus configuration and defer a clean-up
	defer oc.Setup(qr.ServiceName).Close()

	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", qr.ServiceName,
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

	// Create our QR Service
	var svc qr.Service
	{
		svc = implementation.NewService(logger)
		// add service level middlewares here
	}

	// Create our Go kit endpoints for the QR Service
	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(svc)
		// trace our server side endpoints
		endpoints = transport.Endpoints{
			Generate: oc.ServerEndpoint("Generate")(endpoints.Generate),
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
		// set-up our grpc transport
		var (
			bindIP, _     = network.HostIP()
			ocTracing     = kitoc.GRPCServerTrace()
			serverOptions = []kitgrpc.ServerOption{ocTracing}
			qrService     = grpctransport.NewGRPCServer(endpoints, serverOptions, logger)
			listener, _   = net.Listen("tcp", bindIP+":0") // dynamic port assignment
			svcInstance   = fmt.Sprintf("/services/%s/grpc/%s/", qr.ServiceName, instance)
			addr          = listener.Addr().String()
			ttl           = etcd.NewTTLOption(3*time.Second, 10*time.Second)
			serviceEntry  = etcd.Service{Key: svcInstance, Value: addr, TTL: ttl}
			registrar     = etcd.NewRegistrar(sdc, serviceEntry, logger)
			grpcServer    = grpc.NewServer()
		)
		pb.RegisterQRServer(grpcServer, qrService)

		g.Add(func() error {
			registrar.Register()
			return grpcServer.Serve(listener)
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
