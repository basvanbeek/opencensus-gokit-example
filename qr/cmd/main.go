package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/basvanbeek/opencensus-gokit-example"
	"github.com/basvanbeek/opencensus-gokit-example/qr"
	"github.com/basvanbeek/opencensus-gokit-example/qr/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport/grpc"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport/grpc/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger,
			"svc", "QRGenerator",
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var svc qr.Service
	svc = implementation.NewService(logger)

	var endpoints transport.Endpoints
	endpoints = transport.MakeEndpoints(svc)

	var g run.Group
	{
		var (
			grpcQRServer = svcgrpc.NewGRPCServer(endpoints, logger)
			ln, _        = net.Listen("tcp", ocgokitexample.QRAddr)
		)

		g.Add(func() error {
			grpcServer := grpc.NewServer()
			pb.RegisterQRServer(grpcServer, grpcQRServer)
			return grpcServer.Serve(ln)
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
