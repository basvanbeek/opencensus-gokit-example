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

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/oklog/run"

	// project
	"github.com/basvanbeek/opencensus-gokit-example"
	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/frontend/implementation"
	devclient "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/clients/device"
	qrclient "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/clients/grpc/qr"
	svchttp "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/http"
)

func main() {
	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "Frontend",
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
		var err error
		// create our Go kit etcd client
		sdc, err = etcd.NewClient(ctx, []string{"http://localhost:2379"}, etcd.ClientOptions{})
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var svc frontend.Service
	{
		// create an instancer for the device client
		devInstancer, err := etcd.NewInstancer(sdc, "/services/Device/http", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}

		// initialize our Device client using http transport
		devClient := devclient.NewHTTP(devInstancer, logger)

		// create an instancer for the QR client
		qrInstancer, err := etcd.NewInstancer(sdc, "/services/QR/grpc", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}
		// initialize QR client
		qrClient := qrclient.New(qrInstancer, logger)

		// create our frontend service
		svc = implementation.NewService(devClient, qrClient, logger)
		// add service level middlewares here
	}

	var endpoints implementation.Endpoints
	{
		endpoints = implementation.MakeEndpoints(svc)
		// add endpoint level middlewares here
	}

	// run.Group manages our goroutine lifecycles
	// see: https://www.youtube.com/watch?v=LHe1Cb_Ud_M&t=15m45s
	var g run.Group
	{
		// set-up our http transport
		var (
			feService   = svchttp.NewHTTPHandler(endpoints)
			listener, _ = net.Listen("tcp", ocgokitexample.FrontendAddr)
		)

		g.Add(func() error {
			return http.Serve(listener, feService)
		}, func(error) {
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
