package device

import (
	// stdlib
	"context"
	"time"

	// external
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	// project
	devgrpc "github.com/basvanbeek/opencensus-gokit-example/clients/device/grpc"
	devhttp "github.com/basvanbeek/opencensus-gokit-example/clients/device/http"
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/http/routes"
	"github.com/basvanbeek/opencensus-gokit-example/shared/grpcconn"
)

type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

// NewHTTP returns a new device client using the HTTP transport
func NewHTTP(instancer sd.Instancer, logger log.Logger) device.Service {
	// initialize our codec context
	codec := devhttp.Codec{Route: routes.InitEndpoints(mux.NewRouter())}

	// set-up our http transport options
	options := []kithttp.ClientOption{
		kitoc.HTTPClientTrace(),
	}

	// configure rate limiter
	rl := rate.NewLimiter(rate.Every(time.Second), 1000)

	// configure circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Device/http",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})

	// initialize our endpoint middleware
	mw := middleware(rl, cb)

	// initialize our factory handler
	handler := factory(instancer, logger)

	return &client{
		endpoints: transport.Endpoints{
			Unlock: handler(devhttp.NewFactory(codec.Unlock, mw, options...)),
		},
		logger: logger,
	}
}

// NewGRPC returns a new device client using the gRPC transport
func NewGRPC(instancer sd.Instancer, logger log.Logger) device.Service {
	// set-up our grpc transport options
	options := []kitgrpc.ClientOption{
		kitoc.GRPCClientTrace(),
	}

	// initialize our gRPC host mapper helper
	hm := grpcconn.NewHostMapper(grpc.WithInsecure())

	// configure rate limiter
	rl := rate.NewLimiter(rate.Every(time.Second), 1000)

	// configure circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Device/grpc",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})

	// initialize our endpoint middleware
	mw := middleware(rl, cb)

	// initialize our factory handler
	handler := factory(instancer, logger)

	return &client{
		endpoints: transport.Endpoints{
			Unlock: handler(devgrpc.NewFactory(devgrpc.Unlock, mw, hm, options...)),
		},
		logger: logger,
	}

}

func (c client) Unlock(ctx context.Context, eventID, deviceID uuid.UUID, code string) (*device.Session, error) {
	res, err := c.endpoints.Unlock(ctx, transport.UnlockRequest{
		EventID:  eventID,
		DeviceID: deviceID,
		Code:     code,
	})
	if err != nil {
		return nil, err
	}
	response := res.(transport.UnlockResponse)
	return &device.Session{
		EventCaption:  response.EventCaption,
		DeviceCaption: response.DeviceCaption,
	}, nil
}

// middleware wraps a client endpoint with middlewares
func middleware(rl *rate.Limiter, cb *gobreaker.CircuitBreaker) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		e = ratelimit.NewErroringLimiter(rl)(e)
		e = circuitbreaker.Gobreaker(cb)(e)
		return e
	}
}

// factory creates a service discovery driven Go kit client endpoint
func factory(instancer sd.Instancer, logger log.Logger) func(sd.Factory) endpoint.Endpoint {
	return func(factory sd.Factory) endpoint.Endpoint {
		// endpointer manages list of available endpoints servicing our method
		endpointer := sd.NewEndpointer(instancer, factory, logger)

		// balancer can do a round robin pick from the endpointer list
		balancer := lb.NewRoundRobin(endpointer)

		// retry uses balancer for executing a method call with retry and timeout
		// logic so client consumer does not have to think about it.
		return lb.Retry(3, 5*time.Second, balancer)
	}
}
