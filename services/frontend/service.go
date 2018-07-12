package frontend

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/kevinburke/go.uuid"
)

// ServiceName of this service.
const ServiceName = "frontend"

// Service describes our Frontend service.
type Service interface {
	Login(ctx context.Context, user, pass string) (*Login, error)

	EventCreate(ctx context.Context, tenantID uuid.UUID, event Event) (*uuid.UUID, error)
	EventGet(ctx context.Context, tenantID, eventID uuid.UUID) (*Event, error)
	EventUpdate(ctx context.Context, tenantID uuid.UUID, event Event) error
	EventDelete(ctx context.Context, tenantID, eventID uuid.UUID) error
	EventList(ctx context.Context, tenantID uuid.UUID) ([]*Event, error)

	UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (*Session, error)

	GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error)
}

// Common Service Errors
var (
	ErrUserPassRequired  = errors.New("both user and pass are required")
	ErrUserPassUnknown   = errors.New("unknown user/pass combination")
	ErrRequireEventID    = errors.New("missing required event id")
	ErrRequireDeviceID   = errors.New("missing required device id")
	ErrRequireUnlockCode = errors.New("missing required unlock code")
	ErrEventNotFound     = errors.New("event not found")
	ErrUnlockNotFound    = errors.New("device / unlock code combination not found")
)

// Login holds login details
type Login struct {
	ID         uuid.UUID
	Name       string
	TenantID   uuid.UUID
	TenantName string
}

// Event holds event details
type Event struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Session holds session details
type Session struct {
	EventID       uuid.UUID `json:"event_id,omitempty"`
	EventCaption  string    `json:"event_caption,omitempty"`
	DeviceID      uuid.UUID `json:"device_id,omitempty"`
	DeviceCaption string    `json:"device_caption,omitempty"`
	Token         string    `json:"token,omitempty"`
}
