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
	"google.golang.org/grpc"

	"github.com/basvanbeek/opencensus-gokit-example"
	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/frontend/implementation"
	qrclient "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/clients/qr"
	svchttp "github.com/basvanbeek/opencensus-gokit-example/frontend/transport/http"
	"github.com/basvanbeek/opencensus-gokit-example/qr"
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

	var qrClient qr.Service
	{
		conn, _ := grpc.Dial(ocgokitexample.QRAddr, grpc.WithInsecure())
		qrClient = qrclient.New(conn, logger)
	}

	var svc frontend.Service
	svc = implementation.NewService(qrClient, logger)

	var endpoints implementation.Endpoints
	endpoints = implementation.MakeEndpoints(svc)

	var g run.Group
	{
		var (
			hnd   = svchttp.NewHTTPHandler(endpoints)
			ln, _ = net.Listen("tcp", ocgokitexample.FrontendAddr)
		)

		g.Add(func() error {
			return http.Serve(ln, hnd)
		}, func(error) {
			ln.Close()
		})
	}
	{
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

	level.Error(logger).Log("exit", g.Run())
}
