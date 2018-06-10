package svcgrpc

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	uuid "github.com/satori/go.uuid"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/transport/pb"
)

// grpc transport service for QR service.
type grpcServer struct {
	unlock grpctransport.Handler
	logger log.Logger
}

// NewGRPCServer returns a new gRPC service for the provided Go kit endpoints
func NewGRPCServer(
	endpoints transport.Endpoints, options []grpctransport.ServerOption,
	logger log.Logger,
) pb.DeviceServer {
	var (
		errorLogger = grpctransport.ServerErrorLogger(logger)
	)

	options = append(options, errorLogger)

	return &grpcServer{
		unlock: grpctransport.NewServer(
			endpoints.Unlock, decodeUnlockRequest, encodeUnlockResponse, options...,
		),
		logger: logger,
	}
}

// Generate glues the gRPC method to the Go kit service method
func (s *grpcServer) Unlock(ctx oldcontext.Context, req *pb.UnlockRequest) (*pb.UnlockResponse, error) {
	_, rep, err := s.unlock.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UnlockResponse), nil
}

// decodeUnlockRequest decodes the incoming grpc payload to our go kit payload
func decodeUnlockRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UnlockRequest)
	return transport.UnlockRequest{
		EventID:  uuid.FromBytesOrNil(req.EventId),
		DeviceID: uuid.FromBytesOrNil(req.DeviceId),
		Code:     req.Code,
	}, nil
}

// encodeUnlockResponse encodes the outgoing go kit payload to the grpc payload
func encodeUnlockResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.UnlockResponse)
	switch res.Err {
	case nil:
		return &pb.UnlockResponse{
			EventCaption:  res.EventCaption,
			DeviceCaption: res.DeviceCaption,
		}, nil
	case device.ErrRequireEventID, device.ErrRequireDeviceID, device.ErrRequireUnlockCode:
		return nil, status.Error(codes.InvalidArgument, res.Err.Error())
	case device.ErrEventNotFound, device.ErrUnlockNotFound:
		return nil, status.Error(codes.Unauthenticated, res.Err.Error())
	default:
		return nil, status.Error(codes.Unknown, res.Err.Error())
	}
}
