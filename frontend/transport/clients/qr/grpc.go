package qr

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/basvanbeek/opencensus-gokit-example/qr"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport/grpc/pb"
)

type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

// New returns a new QR client using gRPC transport
func New(conn *grpc.ClientConn, logger log.Logger) qr.Service {
	// configure circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})
	// configure rate limiter
	rl := rate.NewLimiter(rate.Every(time.Second), 100)

	options := []grpctransport.ClientOption{}

	// grpc client transport endpoint
	var generateEndpoint endpoint.Endpoint
	{
		generateClient := grpctransport.NewClient(
			conn, "pb.QR", "Generate",
			encodeGenerateRequest, decodeGenerateResponse, pb.GenerateResponse{},
			options...,
		)
		// endpoint middlewares
		generateEndpoint = generateClient.Endpoint()
		generateEndpoint = ratelimit.NewErroringLimiter(rl)(generateEndpoint)
		generateEndpoint = circuitbreaker.Gobreaker(cb)(generateEndpoint)
	}

	return client{
		endpoints: transport.Endpoints{
			Generate: generateEndpoint,
		},
		logger: logger,
	}
}

func (c client) Generate(
	ctx context.Context, data string, recLevel qr.RecoveryLevel, size int,
) ([]byte, error) {
	// we can also validate parameters before sending the request
	if recLevel < qr.LevelL || recLevel > qr.LevelH {
		return nil, qr.ErrInvalidRecoveryLevel
	}
	if size > 4096 {
		return nil, qr.ErrInvalidSize
	}

	// call our client side go kit endpoint
	res, err := c.endpoints.Generate(
		ctx,
		transport.GenerateRequest{Data: data, Level: recLevel, Size: size},
	)
	gErr := status.Convert(err)
	switch gErr.Code() {
	case codes.Unknown:
		return nil, errors.New(gErr.Message())
	case codes.InvalidArgument:
		return nil, errors.New(gErr.Message())
	}
	response := res.(transport.GenerateResponse)
	return response.QR, nil
}

func encodeGenerateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(transport.GenerateRequest)
	return &pb.GenerateRequest{
		Data:  req.Data,
		Level: int32(req.Level),
		Size:  int32(req.Size),
	}, nil
}

func decodeGenerateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GenerateResponse)
	return transport.GenerateResponse{QR: resp.Image}, nil
}
