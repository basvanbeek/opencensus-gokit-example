package transport

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/endpoint"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
)

// Endpoints holds all Go kit endpoints for the service.
type Endpoints struct {
	Generate endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the service.
func MakeEndpoints(s qr.Service) Endpoints {
	return Endpoints{
		Generate: makeGenerateEndpoint(s),
	}
}

func makeGenerateEndpoint(s qr.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenerateRequest)
		qr, err := s.Generate(ctx, req.Data, req.Level, req.Size)
		return GenerateResponse{QR: qr, Err: err}, nil
	}
}
