package transport

import (
	"context"

	"github.com/basvanbeek/opencensus-gokit-example/qr"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Generate endpoint.Endpoint
}

func MakeEndpoints(s qr.Service) Endpoints {
	return Endpoints{
		Generate: makeGenerateEndpoint(s),
	}
}

func makeGenerateEndpoint(s qr.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenerateRequest)
		qr, err := s.Generate(ctx, req.Data, req.Level, req.Size)
		return GenerateResponse{QR: qr, err: err}, nil
	}
}
