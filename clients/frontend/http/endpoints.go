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
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend/transport/http/routes"
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
		Login: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"Login",
			factory.EncodeGenericRequest(route.Login),
			decodeLoginResponse,
		),
		EventCreate: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"EventCreate",
			factory.EncodeGenericRequest(route.EventCreate),
			decodeEventCreateResponse,
		),
		EventGet: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"EventGet",
			factory.EncodeGenericRequest(route.EventGet),
			decodeEventGetResponse,
		),
		EventUpdate: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"EventUpdate",
			factory.EncodeGenericRequest(route.EventUpdate),
			decodeEventUpdateResponse,
		),
		EventDelete: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"EventDelete",
			factory.EncodeGenericRequest(route.EventDelete),
			decodeEventDeleteResponse,
		),
		EventList: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"EventList",
			factory.EncodeGenericRequest(route.EventList),
			decodeEventListResponse,
		),
		UnlockDevice: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"UnlockDevice",
			factory.EncodeGenericRequest(route.UnlockDevice),
			decodeUnlockDeviceResponse,
		),
		GenerateQR: factory.CreateHTTPEndpoint(
			instancer,
			middlewares,
			"GenerateQR",
			factory.EncodeGenericRequest(route.GenerateQR),
			decodeGenerateQRResponse,
		),
	}
}
