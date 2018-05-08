package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"

	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/frontend/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/frontend/transport/svchttp"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger,
			"svc", "frontend",
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var svc frontend.Service
	svc = implementation.NewService()

	var endpoints implementation.Endpoints
	endpoints = implementation.MakeEndpoints(svc)

	var g run.Group
	{
		hnd := svchttp.NewHTTPHandler(endpoints)
		ln, _ := net.Listen("tcp", ":8000")
		g.Add(func() error {
			return http.Serve(ln, hnd)
		}, func(error) {
			ln.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 2)
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

	level.Error(logger).Log("exit", g.Run())
}
