package database

import (
	// stdlib
	"context"

	// external
	"github.com/satori/go.uuid"
)

// Repository describes the resource methods needed for this service.
type Repository interface {
	GetDevice(ctx context.Context, eventID, deviceID uuid.UUID) (*Session, error)
}

// Session holds session details
type Session struct {
	EventCaption  string
	DeviceCaption string
}
