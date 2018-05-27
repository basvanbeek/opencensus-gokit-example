package grpc

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/services/qr/transport/pb"
)

// grpc transport service for QR service.
type grpcServer struct {
	generate grpctransport.Handler
	logger   log.Logger
}

// NewGRPCServer returns a new gRPC service for the provided Go kit endpoints
func NewGRPCServer(endpoints transport.Endpoints, logger log.Logger) pb.QRServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		generate: grpctransport.NewServer(
			endpoints.Generate, decodeGenerateRequest, encodeGenerateResponse, options...,
		),
		logger: logger,
	}
}

// Generate glues the gRPC method to the Go kit service method
func (s *grpcServer) Generate(ctx oldcontext.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	_, rep, err := s.generate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GenerateResponse), nil
}

// decodeGenerateRequest decodes the incoming grpc payload to our go kit payload
func decodeGenerateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GenerateRequest)
	return transport.GenerateRequest{
		Data:  req.Data,
		Level: qr.RecoveryLevel(req.Level),
		Size:  int(req.Size),
	}, nil
}

// encodeGenerateResponse encodes the outgoing go kit payload to the grpc payload
func encodeGenerateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.GenerateResponse)
	switch res.Err {
	case nil:
		return &pb.GenerateResponse{Image: res.QR}, nil
	case qr.ErrInvalidRecoveryLevel, qr.ErrInvalidSize:
		return nil, status.Error(codes.InvalidArgument, res.Err.Error())
	default:
		return nil, status.Error(codes.Unknown, res.Err.Error())
	}
}
