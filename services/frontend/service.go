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

// Frontend Service Error descriptions
const (
	ErrorService           = "internal service error"
	ErrorUnauthorized      = "unauthorized"
	ErrorUserPassRequired  = "both user and pass are required"
	ErrorUserPassUnknown   = "unknown user/pass combination"
	ErrorRequireEventID    = "missing required event id"
	ErrorRequireDeviceID   = "missing required device id"
	ErrorRequireUnlockCode = "missing required unlock code"
	ErrorEventNotFound     = "event not found"
	ErrorEventExists       = "event already exists"
	ErrorUnlockNotFound    = "device / unlock code combination not found"
	ErrorInvalidQRParams   = "QR Code can't be generated using provided parameters"
	ErrorQRGenerate        = "QR Code generator failed"
)

// Frontend Service Errors
var (
	ErrService           = errors.New(ErrorService)
	ErrUnauthorized      = errors.New(ErrorUnauthorized)
	ErrUserPassRequired  = errors.New(ErrorUserPassRequired)
	ErrUserPassUnknown   = errors.New(ErrorUserPassUnknown)
	ErrRequireEventID    = errors.New(ErrorRequireEventID)
	ErrRequireDeviceID   = errors.New(ErrorRequireDeviceID)
	ErrRequireUnlockCode = errors.New(ErrorRequireUnlockCode)
	ErrEventNotFound     = errors.New(ErrorEventNotFound)
	ErrEventExists       = errors.New(ErrorEventExists)
	ErrUnlockNotFound    = errors.New(ErrorUnlockNotFound)

	ErrInvalidQRParams = errors.New(ErrorInvalidQRParams)
	ErrQRGenerate      = errors.New(ErrorQRGenerate)
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
