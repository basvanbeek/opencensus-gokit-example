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
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/satori/go.uuid"
	"go.opencensus.io/trace"

	// project
	feclient "github.com/basvanbeek/opencensus-gokit-example/clients/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/shared/oc"
)

const (
	serviceName = "cli"
)

func main() {
	var (
		err      error
		instance = uuid.Must(uuid.NewV4())
	)

	// initialize our OpenCensus configuration and defer a clean-up
	defer oc.Setup(serviceName).Close()

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
		var instancer sd.Instancer
		instancer, err = etcd.NewInstancer(sdc, "/services/"+frontend.ServiceName+"/http", logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
		}

		client = feclient.NewHTTPClient(instancer, logger)
	}

	var tenantID uuid.UUID
	{
		ctx, span := trace.StartSpan(ctx, "Do Login")
		details, err := client.Login(ctx, "john", "doe")
		if err != nil {
			span.SetStatus(trace.Status{
				Code:    trace.StatusCodeUnknown,
				Message: err.Error(),
			})
			span.End()
			level.Error(logger).Log("msg", "login failed", "exit", err)
			os.Exit(-1)
		}
		span.SetStatus(trace.Status{Code: trace.StatusCodeOK})
		span.End()
		level.Debug(logger).Log("msg", "login succeeded", "details", fmt.Sprintf("%+v", details))
		tenantID = details.TenantID
	}

	{
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		ctx, span := trace.StartSpan(ctx, "Do EventCreate")
		id, err := client.EventCreate(ctx, tenantID, frontend.Event{
			Name: "Marine Corps Marathon",
		})
		if err != nil {
			span.SetStatus(trace.Status{
				Code:    trace.StatusCodeUnknown,
				Message: err.Error(),
			})
			level.Error(logger).Log("msg", "event create failed", "exit", err)
		} else {
			span.SetStatus(trace.Status{Code: trace.StatusCodeOK})
			level.Debug(logger).Log("msg", "event create succeeded", "id", id.String())
		}
		span.End()
	}

	{
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		ctx, span := trace.StartSpan(ctx, "Do EventList")
		events, err := client.EventList(ctx, tenantID)

		// highlight how terribly expensive a call to spew is...
		span.Annotate(nil, "spew.dump:start")
		fmt.Printf("\nCLIENT EVENT LIST:\nRES:%s\nERR: %+v\n\n", spew.Sdump(events), err)
		span.Annotate(nil, "spew.dump:end")
		if err != nil {
			span.SetStatus(trace.Status{
				Code:    trace.StatusCodeUnknown,
				Message: err.Error(),
			})
		} else {
			span.SetStatus(trace.Status{Code: trace.StatusCodeOK})
		}
		span.End()
	}

	// details, err = client.Login(ctx, "jane", "doe")
	// fmt.Printf("\nCLIENT LOGIN:\nRES:%+v\nERR: %+v\n\n", details, err)
	//
	// details, err = client.Login(ctx, "Anonymous", "Coward")
	// fmt.Printf("\nCLIENT LOGIN:\nRES:%+v\nERR: %+v\n\n", details, err)

}
