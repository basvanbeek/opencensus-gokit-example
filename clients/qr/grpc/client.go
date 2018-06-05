package grpcclient

import (
	// stdlib

	// external

	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	// project

	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/grpcconn"
	"github.com/basvanbeek/opencensus-gokit-example/shared/loggermw"
)

// InitEndpoints returns an initialized set of Go kit gRPC endpoints
func InitEndpoints(instancer sd.Instancer, logger log.Logger) transport.Endpoints {
	// initialize our gRPC host mapper helper
	hm := grpcconn.NewHostMapper(grpc.WithInsecure())

	// configure client wide rate limiter for all instances and all method
	// endpoints
	rl := ratelimit.NewErroringLimiter(
		rate.NewLimiter(rate.Every(time.Second), 1000),
	)

	// debug logging middleware
	lmw := loggermw.LoggerMiddleware(level.Debug(logger))

	// chain our service wide middlewares
	middlewares := endpoint.Chain(lmw, rl)

	return transport.Endpoints{
		Generate: createEndpoint(
			instancer,
			hm,
			middlewares,
			"Generate",
			pb.GenerateResponse{},
			encodeGenerateRequest,
			decodeGenerateResponse,
			decodeGenerateError(),
		),
	}
}
