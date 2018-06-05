package grpc

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport/pb"
)

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

func decodeGenerateError() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			// call our gRPC client endpoint
			response, err := e(ctx, request)
			// check error response
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.OK:
				// no error encountered... proceed with regular response payload
				return response, nil
			case codes.InvalidArgument, codes.FailedPrecondition:
				// business logic error which should not be retried or trigger
				// the circuitbreaker.
				return transport.GenerateResponse{Err: errors.New(st.Message())}, nil
			default:
				// error which might invoke a retry or trigger a circuitbreaker
				return nil, errors.New(st.Message())
			}
		}
	}
}
