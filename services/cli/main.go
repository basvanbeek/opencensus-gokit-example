package main

import (
	// stdlib
	"context"
	"fmt"
	"os"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcd"

	// project
	feclient "github.com/basvanbeek/opencensus-gokit-example/services/cli/transport/clients/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
)

func main() {
	// initialize our structured logger for the service
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "Frontend-CLI",
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	var client frontend.Service
	{
		// create an instancer for the event client
		instancer, err := etcd.NewInstancer(sdc, "/services/Frontend/http", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}

		client = feclient.NewHTTP(instancer, logger)
	}

	details, err := client.Login(ctx, "john", "doe")
	fmt.Printf("CLIENT LOGIN: %+v %+v\n", details, err)

}
