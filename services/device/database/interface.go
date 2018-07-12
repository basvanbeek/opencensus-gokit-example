package database

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/kevinburke/go.uuid"
)

// Common Errors
var (
	ErrRepository = errors.New("unable to handle request")
	ErrNotFound   = errors.New("device not found")
)

// Repository describes the resource methods needed for this service.
type Repository interface {
	GetDevice(ctx context.Context, eventID, deviceID uuid.UUID) (*Session, error)
}

// Session holds session details
type Session struct {
	EventCaption  string
	DeviceCaption string
	UnlockHash    []byte
}
