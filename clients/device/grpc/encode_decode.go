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
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/pb"
)

func encodeUnlockRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(transport.UnlockRequest)
	return &pb.UnlockRequest{
		EventId: req.EventID.Bytes(),
	}, nil
}

func decodeUnlockResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*pb.UnlockResponse)
	return transport.UnlockResponse{
		DeviceCaption: res.DeviceCaption,
		EventCaption:  res.EventCaption,
	}, nil
}

func decodeUnlockError() endpoint.Middleware {
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
			case codes.InvalidArgument, codes.Unauthenticated:
				// business logic error which should not retry or trigger
				// the circuitbreaker as service is behaving normally.
				return transport.UnlockResponse{Err: errors.New(st.Message())}, nil
			default:
				// error which might invoke a retry or trigger a circuitbreaker
				return nil, errors.New(st.Message())
			}
		}
	}
}
