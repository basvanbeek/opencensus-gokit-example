package qr

import (
	// stdlib
	"context"
	"errors"
	"io"
	"time"

	// external
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport/pb"
	"github.com/basvanbeek/opencensus-gokit-example/shared/grpcconn"
)

// client grpc transport to QR service.
type client struct {
	endpoints transport.Endpoints
	logger    log.Logger
}

// New returns a new QR client using gRPC transport
func New(instancer sd.Instancer, logger log.Logger) qr.Service {
	// initialize our gRPC host mapper helper
	hm := grpcconn.NewHostMapper(grpc.WithInsecure())

	options := []kitgrpc.ClientOption{}

	// makeEndpoint wires a QR service Go kit method endpoint
	makeEndpoint := func(
		method string, reply interface{},
		enc kitgrpc.EncodeRequestFunc, dec kitgrpc.DecodeResponseFunc,
	) endpoint.Endpoint {

		// configure circuit breaker
		cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "QR/grpc/" + method,
			MaxRequests: 5,
			Interval:    10 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 5
			},
		})

		// configure rate limiter
		rl := rate.NewLimiter(rate.Every(time.Second), 100)

		// our method sd.Factory is called when a new QR service is discovered.
		factory := func(instance string) (endpoint.Endpoint, io.Closer, error) {

			// try to get connection to advertised instance
			conn, closer, err := hm.Get(instance)
			if err != nil {
				return nil, nil, err
			}

			// set-up our go kit client endpoint
			var e endpoint.Endpoint
			e = kitgrpc.NewClient(conn, "pb.QR", method, enc, dec, reply, options...).Endpoint()
			e = ratelimit.NewErroringLimiter(rl)(e)
			e = circuitbreaker.Gobreaker(cb)(e)

			return e, closer, nil
		}

		// endpointer manages list of available endpoints servicing our method
		endpointer := sd.NewEndpointer(instancer, factory, logger)

		// balancer can do a round robin pick from the endpointer list
		balancer := lb.NewRoundRobin(endpointer)

		// retry uses balancer for executing a method call with retry and timeout
		// logic so client consumer does not have to think about it.
		retry := lb.Retry(3, 5*time.Second, balancer)
		return retry
	}

	// create our QR client by initializing all method endpoints
	return client{
		endpoints: transport.Endpoints{
			Generate: makeEndpoint("Generate", pb.GenerateResponse{}, encodeGenerateRequest, decodeGenerateResponse),
		},
		logger: logger,
	}
}

// Generate calls the QR Service Generate method.
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

// encodeGenerateRequest encodes the outgoing go kit payload to the grpc payload
func encodeGenerateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(transport.GenerateRequest)
	return &pb.GenerateRequest{
		Data:  req.Data,
		Level: int32(req.Level),
		Size:  int32(req.Size),
	}, nil
}

// decodeGenerateResponse decodes the incoming grpc payload to go kit payload
func decodeGenerateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GenerateResponse)
	return transport.GenerateResponse{QR: resp.Image}, nil
}
