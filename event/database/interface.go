package database

import (
	// stdlib
	"context"
	"errors"

	// external

	"github.com/satori/go.uuid"
)

// Common Errors
var (
	ErrRepository = errors.New("unable to handle request")
	ErrNotFound   = errors.New("event not found")
	ErrIDExists   = errors.New("event id already exists")
	ErrNameExists = errors.New("event name already exists")
)

// Repository describes the resource methods needed for this service.
type Repository interface {
	Create(ctx context.Context, event Event) (*uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*Event, error)
	Update(ctx context.Context, event Event) error
	Delete(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID) ([]*Event, error)
}

// Event holds event details
type Event struct {
	ID       uuid.UUID
	TenantID uuid.UUID
	Name     string
}
