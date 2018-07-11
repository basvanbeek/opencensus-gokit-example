package http

import (
	// stdlib
	"time"

	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/http/routes"
	"github.com/basvanbeek/opencensus-gokit-example/shared/factory"
	"github.com/basvanbeek/opencensus-gokit-example/shared/loggermw"
)

// InitEndpoints returns an initialized set of Go kit HTTP endpoints.
func InitEndpoints(instancer sd.Instancer, logger log.Logger) transport.Endpoints {
	route := routes.Initialize(mux.NewRouter())

	// configure client wide rate limiter for all instances and all method
	// endpoints
	rl := ratelimit.NewErroringLimiter(
		rate.NewLimiter(rate.Every(time.Second), 1000),
	)

	// debug logging middleware
	lmw := loggermw.LoggerMiddleware(level.Debug(logger))

	// chain our service wide middlewares
	middlewares := endpoint.Chain(lmw, rl)

	// create our client endpoints
	return transport.Endpoints{
		Unlock: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"Unlock",
			encodeUnlockRequest(route.Unlock),
			decodeUnlockResponse,
		),
	}
}
