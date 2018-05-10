package main

import (
	// stdlib
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"

	// project
	"github.com/basvanbeek/opencensus-gokit-example"
	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/frontend/implementation"
	qrclient "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/clients/qr"
	svchttp "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/http"
)

func main() {
	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger,
			"svc", "Frontend",
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var svc frontend.Service
	{
		// initialize QR client
		qrClient := qrclient.New(logger)

		// create our frontend service
		svc = implementation.NewService(qrClient, logger)
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
