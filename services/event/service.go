package event

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/satori/go.uuid"
)

// ServiceName of this service.
const ServiceName = "event"

// Service describes our Event service.
type Service interface {
	Create(ctx context.Context, tenantID uuid.UUID, event Event) (*uuid.UUID, error)
	Get(ctx context.Context, tenantID, id uuid.UUID) (*Event, error)
	Update(ctx context.Context, tenantID uuid.UUID, event Event) error
	Delete(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID) ([]*Event, error)
}

// Event Service Error descriptions
const (
	ErrorService      = "internal service error"
	ErrorUnauthorized = "unauthorized"
	ErrorNotFound     = "event not found"
	ErrorEventExists  = "event already exists"
)

// Event Service Errors
var (
	ErrService      = errors.New(ErrorService)
	ErrUnauthorized = errors.New(ErrorUnauthorized)
	ErrNotFound     = errors.New(ErrorNotFound)
	ErrEventExists  = errors.New(ErrorEventExists)
)

// Event data
type Event struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
