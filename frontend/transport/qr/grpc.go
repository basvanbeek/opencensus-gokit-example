package qr

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
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

type Client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) qr.Service {
	// initialize circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{})
	// initialize rate limiter
	rl := rate.NewLimiter(rate.Every(10*time.Second), 2)

	options := []grpctransport.ClientOption{}

	// grpc client transport endpoint
	e := transport.Endpoints{
		Generate: grpctransport.NewClient(
			conn,
			"pb.QR",
			"Generate",
			encodeGenerateRequest,
			decodeGenerateResponse,
			pb.GenerateResponse{},
			options...,
		).Endpoint(),
	}

	// endpoint middlewares
	e.Generate = ratelimit.NewErroringLimiter(rl)(e.Generate)
	e.Generate = circuitbreaker.Gobreaker(cb)(e.Generate)

	return Client{
		endpoints: e,
		logger:    logger,
	}
}

func (c Client) Generate(
	ctx context.Context, data string, recLevel qr.RecoveryLevel, size int,
) ([]byte, error) {
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
