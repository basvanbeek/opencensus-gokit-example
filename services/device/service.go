package device

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/satori/go.uuid"
)

// ServiceName of this service.
const ServiceName = "device"

// Service describes our Device service.
type Service interface {
	Unlock(ctx context.Context, eventID, deviceID uuid.UUID, code string) (*Session, error)
}

// Device Service Error descriptions
const (
	ErrorRequireEventID    = "missing required event id"
	ErrorRequireDeviceID   = "missing required device id"
	ErrorRequireUnlockCode = "missing required unlock code"
	ErrorRepository        = "unable to query repository"
	ErrorEventNotFound     = "event not found"
	ErrorUnlockNotFound    = "device / unlock code combination not found"
)

// Device Service Errors
var (
	ErrRequireEventID    = errors.New(ErrorRequireEventID)
	ErrRequireDeviceID   = errors.New(ErrorRequireDeviceID)
	ErrRequireUnlockCode = errors.New(ErrorRequireUnlockCode)
	ErrRepository        = errors.New(ErrorRepository)
	ErrEventNotFound     = errors.New(ErrorEventNotFound)
	ErrUnlockNotFound    = errors.New(ErrorUnlockNotFound)
)

// Session holds session details
type Session struct {
	EventCaption  string `json:"event_caption,omitempty"`
	DeviceCaption string `json:"device_caption,omitempty"`
	Token         string `json:"token,omitempty"`
}
