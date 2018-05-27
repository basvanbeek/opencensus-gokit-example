package implementation

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/endpoint"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/frontend"
)

// Endpoints holds all Go kit endpoints for the service.
type Endpoints struct {
	Login        endpoint.Endpoint
	UnlockDevice endpoint.Endpoint
	GenerateQR   endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the service.
func MakeEndpoints(s frontend.Service) Endpoints {
	return Endpoints{
		Login:        makeLoginEndpoint(s),
		UnlockDevice: makeUnlockDeviceEndpoint(s),
		GenerateQR:   makeGenerateQREndpoint(s),
	}
}

func makeLoginEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		login, err := s.Login(ctx, req.Login, req.Pass)
		if err != nil {
			return LoginResponse{err: err}, nil
		}
		return LoginResponse{
			ID:         login.ID,
			Name:       login.Name,
			TenantID:   login.TenantID,
			TenantName: login.TenantName,
		}, nil
	}
}

func makeUnlockDeviceEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UnlockDeviceRequest)
		session, err := s.UnlockDevice(ctx, req.EventID, req.DeviceID, req.UnlockCode)
		return UnlockDeviceResponse{Session: session, err: err}, nil
	}
}

func makeGenerateQREndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenerateQRRequest)
		qr, err := s.GenerateQR(ctx, req.EventID, req.DeviceID, req.UnlockCode)
		return GenerateQRResponse{QR: qr, err: err}, nil
	}
}
