package httpclient

import (
	// stdlib

	"io"
	"net/url"
	"time"

	// external
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
	"go.opencensus.io/trace"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/shared/oc"
)

// CreateEndpoint creates a Go kit client endpoint
func CreateEndpoint(
	instancer sd.Instancer, middleware endpoint.Middleware, operationName string,
	route *mux.Route, decodeResponse kithttp.DecodeResponseFunc,
) endpoint.Endpoint {
	options := []kithttp.ClientOption{
		kitoc.HTTPClientTrace(), // OpenCensus HTTP Client transport tracing
	}

	// factory is called each time a new instance is received from service
	// discovery. it will create a new Go kit client endpoint which will be
	// consumed by the endpointer logic.
	factory := func(instance string) (endpoint.Endpoint, io.Closer, error) {
		baseURL, err := url.Parse(instance)
		if err != nil {
			// invalid instance string received... can't build endpoint
			return nil, nil, err
		}

		// set-up our Go kit client endpoint
		// method is not set yet as it will be decided by the provided route
		// when encoding the request using our generic request encoder.
		clientEndpoint := kithttp.NewClient(
			"", baseURL, encodeGenericRequest(route), decodeResponse, options...,
		).Endpoint()

		// configure per instance circuit breaker middleware
		cb := circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:        "CLI/http",
				MaxRequests: 5,
				Interval:    10 * time.Second,
				Timeout:     10 * time.Second,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.ConsecutiveFailures > 5
				},
			}),
		)

		// middleware to trace our client endpoint
		tr := oc.ClientEndpoint(operationName)

		// chain our middlewares
		middleware = endpoint.Chain(cb, tr, middleware)

		return middleware(clientEndpoint), nil, nil
	}

	// endpoints manages the list of available endpoints servicing our method
	endpoints := sd.NewEndpointer(instancer, factory, log.NewNopLogger())

	// balancer can do a random pick from the endpoint list
	balancer := lb.NewRandom(endpoints, time.Now().UnixNano())

	var (
		retryCount    = 3
		retryDuration = 5 * time.Second
	)

	// retryTracer instruments our HTTP endpoint retry/load balancer logic
	retryTracer := kitoc.TraceEndpoint(
		"kit/retry "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("kit.balancer.type", "random"),
			trace.StringAttribute("kit.retry.timeout", retryDuration.String()),
			trace.Int64Attribute("kit.retry.count", int64(retryCount)),
		),
	)

	// retry uses balancer for executing a method call with retry and
	// timeout logic so client consumer does not have to think about it.
	endpoint := lb.Retry(retryCount, retryDuration, balancer)

	// wrap our retries in a parent span
	endpoint = retryTracer(endpoint)

	// return our endpoint
	return endpoint
}
