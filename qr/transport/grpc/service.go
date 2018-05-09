package svcgrpc

import (
	"context"

	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/basvanbeek/opencensus-gokit-example/qr"
	"github.com/basvanbeek/opencensus-gokit-example/qr/implementation"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport"
	"github.com/basvanbeek/opencensus-gokit-example/qr/transport/grpc/pb"
)

type grpcServer struct {
	generate grpctransport.Handler
	logger   log.Logger
}

// NewGRPCServer returns a new gRPC service for provided Go kit endpoints
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

func decodeGenerateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GenerateRequest)
	return transport.GenerateRequest{
		Data:  req.Data,
		Level: qr.RecoveryLevel(req.Level),
		Size:  int(req.Size),
	}, nil
}

func encodeGenerateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(transport.GenerateResponse)
	switch res.Failed() {
	case nil:
		return &pb.GenerateResponse{Image: res.QR}, nil
	case implementation.ErrInvalidRecoveryLevel, implementation.ErrInvalidSize:
		return nil, status.Error(codes.InvalidArgument, res.Failed().Error())
	default:
		return nil, status.Error(codes.Unknown, res.Failed().Error())
	}
}
