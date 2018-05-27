package transport

import (
	// external
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
)

// project

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

// UpdateRequest holds the request parameters for the Update method.
type UpdateRequest struct {
	TenantID uuid.UUID
	Event    event.Event
}

// UpdateResponse holds the response values for the Update method.
type UpdateResponse struct {
	Err error
}

// DeleteRequest holds the request parameters for the Delete method.
type DeleteRequest struct {
	TenantID uuid.UUID
	ID       uuid.UUID
}

// DeleteResponse holds the response values for the Delete method.
type DeleteResponse struct {
	Err error
}

// ListRequest holds the request parameters for the List method.
type ListRequest struct {
	TenantID uuid.UUID
}

// ListResponse holds the response values for the List method.
type ListResponse struct {
	Events []*event.Event
	Err    error
}
