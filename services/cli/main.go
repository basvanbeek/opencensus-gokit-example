package main

import (
	// stdlib
	"context"
	"fmt"
	"os"

	// external
	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd/etcd"
	uuid "github.com/satori/go.uuid"

	// project
	feclient "github.com/basvanbeek/opencensus-gokit-example/services/cli/transport/clients/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
)

const (
	serviceName = "frontend-cli"
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
			"svc", serviceName,
			"instance", instance,
			"ts", log.DefaultTimestampUTC,
			"clr", log.DefaultCaller,
		)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var sdc etcd.Client
	{
		// create our Go kit etcd client
		sdc, err = etcd.NewClient(ctx, []string{"http://localhost:2379"}, etcd.ClientOptions{})
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var client frontend.Service
	{
		// create an instancer for the frontend client
		instancer, err := etcd.NewInstancer(sdc, "/services/"+frontend.ServiceName+"/http", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}

		client = feclient.NewHTTP(instancer, logger)
	}

	details, err := client.Login(ctx, "john", "doe")
	if err != nil {
		level.Error(logger).Log("msg", "login failed", "exit", err)
		os.Exit(-1)
	}
	level.Debug(logger).Log("msg", "login succeeded", "details", fmt.Sprintf("%+v", details))

	id, err := client.EventCreate(ctx, details.TenantID, frontend.Event{
		Name: "Marine Corps Marathon",
	})
	if err != nil {
		level.Error(logger).Log("msg", "event create failed", "exit", err)
	} else {
		level.Debug(logger).Log("msg", "event create succeeded", "id", id.String())
	}

	events, err := client.EventList(ctx, details.TenantID)
	fmt.Printf("\nCLIENT EVENT LIST:\nRES:%s\nERR: %+v\n\n", spew.Sdump(events), err)

	// details, err = client.Login(ctx, "jane", "doe")
	// fmt.Printf("\nCLIENT LOGIN:\nRES:%+v\nERR: %+v\n\n", details, err)
	//
	// details, err = client.Login(ctx, "Anonymous", "Coward")
	// fmt.Printf("\nCLIENT LOGIN:\nRES:%+v\nERR: %+v\n\n", details, err)

}
