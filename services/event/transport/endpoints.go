package transport

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/endpoint"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
)

// Endpoints holds all Go kit endpoints for the service.
type Endpoints struct {
	Create endpoint.Endpoint
	Get    endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
	List   endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the service.
func MakeEndpoints(s event.Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
		List:   makeListEndpoint(s),
	}
}

func makeCreateEndpoint(s event.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.TenantID, req.Event)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetEndpoint(s event.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		event, err := s.Get(ctx, req.TenantID, req.ID)
		return GetResponse{Event: event, Err: err}, nil
	}
}

func makeUpdateEndpoint(s event.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := s.Update(ctx, req.TenantID, req.Event)
		return UpdateResponse{Err: err}, nil
	}
}

func makeDeleteEndpoint(s event.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.TenantID, req.ID)
		return DeleteResponse{Err: err}, nil
	}
}

func makeListEndpoint(s event.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		events, err := s.List(ctx, req.TenantID)
		return ListResponse{Events: events, Err: err}, nil
	}
}
