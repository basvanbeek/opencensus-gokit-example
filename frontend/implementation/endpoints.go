package implementation

import (
	"context"

	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	UnlockDevice endpoint.Endpoint
	GenerateQR   endpoint.Endpoint
}

func MakeEndpoints(s frontend.Service) Endpoints {
	return Endpoints{
		UnlockDevice: makeUnlockDeviceEndpoint(s),
		GenerateQR:   makeGenerateQREndpoint(s),
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
