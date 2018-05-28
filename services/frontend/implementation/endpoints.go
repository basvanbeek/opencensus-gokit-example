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
	EventCreate  endpoint.Endpoint
	EventGet     endpoint.Endpoint
	EventUpdate  endpoint.Endpoint
	EventDelete  endpoint.Endpoint
	EventList    endpoint.Endpoint
	UnlockDevice endpoint.Endpoint
	GenerateQR   endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the service.
func MakeEndpoints(s frontend.Service) Endpoints {
	return Endpoints{
		Login:        makeLoginEndpoint(s),
		EventCreate:  makeEventCreateEndpoint(s),
		EventGet:     makeEventGetEndpoint(s),
		EventUpdate:  makeEventUpdateEndpoint(s),
		EventDelete:  makeEventDeleteEndpoint(s),
		EventList:    makeEventListEndpoint(s),
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

func makeEventCreateEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EventCreateRequest)
		eventID, err := s.EventCreate(ctx, req.TenantID, req.Event)
		return EventCreateResponse{EventID: eventID, err: err}, nil
	}
}

func makeEventGetEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EventGetRequest)
		event, err := s.EventGet(ctx, req.TenantID, req.EventID)
		return EventGetResponse{Event: event, err: err}, nil
	}
}

func makeEventUpdateEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EventUpdateRequest)
		err := s.EventUpdate(ctx, req.TenantID, req.Event)
		return EventUpdateResponse{err: err}, nil
	}
}

func makeEventDeleteEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EventDeleteRequest)
		err := s.EventDelete(ctx, req.TenantID, req.EventID)
		return EventDeleteResponse{err: err}, nil
	}
}

func makeEventListEndpoint(s frontend.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EventListRequest)
		events, err := s.EventList(ctx, req.TenantID)
		return EventListResponse{Events: events, err: err}, nil
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
