package factory

import (
	// stdlib
	"io"
	"time"

	// external
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/shared/grpcconn"
	"github.com/basvanbeek/opencensus-gokit-example/shared/oc"
)

// CreateGRPCEndpoint wires a QR service Go kit method endpoint
func CreateGRPCEndpoint(
	instancer sd.Instancer, hm grpcconn.HostMapper,
	middleware endpoint.Middleware, method string, reply interface{},
	enc kitgrpc.EncodeRequestFunc, dec kitgrpc.DecodeResponseFunc,
	decError endpoint.Middleware,
) endpoint.Endpoint {
	// Set our Go kit gRPC client options
	options := []kitgrpc.ClientOption{
		kitoc.GRPCClientTrace(), // OpenCensus Go kit gRPC client tracing
	}

	// our method sd.Factory is called when a new QR service is discovered.
	factory := func(instance string) (endpoint.Endpoint, io.Closer, error) {
		// try to get connection to advertised instance
		conn, closer, err := hm.Get(instance)
		if err != nil {
			// unable to get a connection to instance... can't build endpoint
			return nil, nil, err
		}

		// set-up our Go kit client endpoint
		clientEndpoint := kitgrpc.NewClient(
			conn, "pb.QR", method, enc, dec, reply, options...,
		).Endpoint()

		if decError != nil {
			// we have a custom gRPC error router
			clientEndpoint = decError(clientEndpoint)
		}

		// configure circuit breaker
		cb := circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:        "QR/grpc/" + method,
				MaxRequests: 5,
				Interval:    10 * time.Second,
				Timeout:     10 * time.Second,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.ConsecutiveFailures > 5
				},
			}),
		)

		// middleware to trace our client endpoint
		tr := oc.ClientEndpoint(method)

		// chain our middlewares
		middleware = endpoint.Chain(cb, tr, middleware)

		return middleware(clientEndpoint), closer, nil
	}

	// endpointer manages list of available endpoints servicing our method
	endpointer := sd.NewEndpointer(instancer, factory, log.NewNopLogger())

	// balancer can do a round robin pick from the endpointer list
	balancer := lb.NewRoundRobin(endpointer)

	// retry uses balancer for executing a method call with retry and timeout
	// logic so client consumer does not have to think about it.
	var (
		count   = 3
		timeout = 5 * time.Second
	)

	// retry uses balancer for executing a method call with retry and
	// timeout logic so client consumer does not have to think about it.
	endpoint := lb.Retry(count, timeout, balancer)

	// wrap our retries in an annotated parent span
	endpoint = oc.RetryEndpoint(method, oc.RoundRobin, count, timeout)(endpoint)

	// return our endpoint
	return endpoint
}
