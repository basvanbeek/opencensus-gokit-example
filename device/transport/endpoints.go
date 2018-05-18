package transport

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/endpoint"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/device"
)

// Endpoints holds all Go kit endpoints for the service.
type Endpoints struct {
	Unlock endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the service.
func MakeEndpoints(s device.Service) Endpoints {
	return Endpoints{
		Unlock: makeUnlockEndpoint(s),
	}
}

func makeUnlockEndpoint(s device.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UnlockRequest)
		res, err := s.Unlock(ctx, req.EventID, req.DeviceID, req.Code)
		if err != nil {
			return nil, err
		}
		return UnlockResponse{EventCaption: res.DeviceCaption, DeviceCaption: res.DeviceCaption}, nil
	}
}
