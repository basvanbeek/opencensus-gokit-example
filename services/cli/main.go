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
	"go.opencensus.io/trace"

	// project
	feclient "github.com/basvanbeek/opencensus-gokit-example/services/cli/transport/clients/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/shared/opencensus"
)

const (
	serviceName = "cli"
)

func main() {
	var (
		err      error
		instance = uuid.Must(uuid.NewV4())
	)

	// initialize our OpenCensus configuration
	defer opencensus.Setup(serviceName).Close()

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

	var tenantID uuid.UUID
	{
		ctx, span := trace.StartSpan(ctx, "LoginAppSpan")
		details, err := client.Login(ctx, "john", "doe")
		if err != nil {
			span.SetStatus(trace.Status{trace.StatusCodeUnknown, err.Error()})
			span.End()
			level.Error(logger).Log("msg", "login failed", "exit", err)
			os.Exit(-1)
		}
		span.SetStatus(trace.Status{trace.StatusCodeOK, ""})
		span.End()
		level.Debug(logger).Log("msg", "login succeeded", "details", fmt.Sprintf("%+v", details))
		tenantID = details.TenantID
	}

	// id, err := client.EventCreate(ctx, tenantID, frontend.Event{
	// 	Name: "Marine Corps Marathon",
	// })
	// if err != nil {
	// 	level.Error(logger).Log("msg", "event create failed", "exit", err)
	// } else {
	// 	level.Debug(logger).Log("msg", "event create succeeded", "id", id.String())
	// }

	{
		ctx, span := trace.StartSpan(ctx, "EventList")
		events, err := client.EventList(ctx, tenantID)
		fmt.Printf("\nCLIENT EVENT LIST:\nRES:%s\nERR: %+v\n\n", spew.Sdump(events), err)
		if err != nil {
			span.SetStatus(trace.Status{trace.StatusCodeUnknown, err.Error()})
		} else {
			span.SetStatus(trace.Status{trace.StatusCodeOK, ""})
		}
		span.End()
	}

	// details, err = client.Login(ctx, "jane", "doe")
	// fmt.Printf("\nCLIENT LOGIN:\nRES:%+v\nERR: %+v\n\n", details, err)
	//
	// details, err = client.Login(ctx, "Anonymous", "Coward")
	// fmt.Printf("\nCLIENT LOGIN:\nRES:%+v\nERR: %+v\n\n", details, err)

}
