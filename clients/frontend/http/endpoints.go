package httpclient

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
		Login: createEndpoint(
			instancer,
			middlewares,
			"Login",
			route.Login,
			decodeLoginResponse,
		),
		EventCreate: createEndpoint(
			instancer,
			middlewares,
			"EventCreate",
			route.EventCreate,
			decodeEventCreateResponse,
		),
		EventGet: createEndpoint(
			instancer,
			middlewares,
			"EventGet",
			route.EventGet,
			decodeEventGetResponse,
		),
		EventUpdate: createEndpoint(
			instancer,
			middlewares,
			"EventUpdate",
			route.EventUpdate,
			decodeEventUpdateResponse,
		),
		EventDelete: createEndpoint(
			instancer,
			middlewares,
			"EventDelete",
			route.EventDelete,
			decodeEventDeleteResponse,
		),
		EventList: createEndpoint(
			instancer,
			middlewares,
			"EventList",
			route.EventList,
			decodeEventListResponse,
		),
		UnlockDevice: createEndpoint(
			instancer,
			middlewares,
			"UnlockDevice",
			route.UnlockDevice,
			decodeUnlockDeviceResponse,
		),
		GenerateQR: createEndpoint(
			instancer,
			middlewares,
			"GenerateQR",
			route.GenerateQR,
			decodeGenerateQRResponse,
		),
	}
}
