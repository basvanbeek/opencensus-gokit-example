package transport

import (
	// external
	"github.com/go-kit/kit/endpoint"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
)

var (
	_ endpoint.Failer = CreateResponse{}
	_ endpoint.Failer = GetResponse{}
	_ endpoint.Failer = UpdateResponse{}
	_ endpoint.Failer = DeleteResponse{}
	_ endpoint.Failer = ListResponse{}
)

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	TenantID uuid.UUID
	Event    event.Event
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  *uuid.UUID
	Err error
}

// Failed implements Failer
func (r CreateResponse) Failed() error { return r.Err }

// GetRequest holds the request parameters for the Get method.
type GetRequest struct {
	TenantID uuid.UUID
	ID       uuid.UUID
}

// GetResponse holds the response values for the Get method.
type GetResponse struct {
	Event *event.Event
	Err   error
}

// Failed implements Failer
func (r GetResponse) Failed() error { return r.Err }

// UpdateRequest holds the request parameters for the Update method.
type UpdateRequest struct {
	TenantID uuid.UUID
	Event    event.Event
}

// UpdateResponse holds the response values for the Update method.
type UpdateResponse struct {
	Err error
}

// Failed implements Failer
func (r UpdateResponse) Failed() error { return r.Err }

// DeleteRequest holds the request parameters for the Delete method.
type DeleteRequest struct {
	TenantID uuid.UUID
	ID       uuid.UUID
}

// DeleteResponse holds the response values for the Delete method.
type DeleteResponse struct {
	Err error
}

// Failed implements Failer
func (r DeleteResponse) Failed() error { return r.Err }

// ListRequest holds the request parameters for the List method.
type ListRequest struct {
	TenantID uuid.UUID
}

// ListResponse holds the response values for the List method.
type ListResponse struct {
	Events []*event.Event
	Err    error
}

// Failed implements Failer
func (r ListResponse) Failed() error { return r.Err }
