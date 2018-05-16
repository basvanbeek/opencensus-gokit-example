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
	feimplementation "github.com/basvanbeek/opencensus-gokit-example/frontend/implementation"
	fesvchttp "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/http"
	"github.com/basvanbeek/opencensus-gokit-example/qr"
	qrimplementation "github.com/basvanbeek/opencensus-gokit-example/qr/implementation"
)

func main() {
	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "ElegantMonolith",
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Create our QR service component
	var qrService qr.Service
	{
		qrService = qrimplementation.NewService(log.With(logger, "component", "QR"))
	}

	// Create our frontend service component
	var frontendService frontend.Service
	{
		frontendService = feimplementation.NewService(qrService, log.With(logger, "component", "QR"))
		// add service level middlewares here
	}

	var endpoints feimplementation.Endpoints
	{
		endpoints = feimplementation.MakeEndpoints(frontendService)
		// add frontend endpoint level middlewares here
	}

	// run.Group manages our goroutine lifecycles
	// see: https://www.youtube.com/watch?v=LHe1Cb_Ud_M&t=15m45s
	var g run.Group
	{
		// set-up our http transport
		var (
			frontendService = fesvchttp.NewHTTPHandler(endpoints)
			listener, _     = net.Listen("tcp", ocgokitexample.FrontendAddr)
		)

		g.Add(func() error {
			return http.Serve(listener, frontendService)
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
