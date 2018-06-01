package device

import (
	// stdlib
	"context"
	"errors"

	// external
	"github.com/satori/go.uuid"
)

// Service name
const (
	ServiceName = "device"
)

// Service describes our Device service.
type Service interface {
	Unlock(ctx context.Context, eventID, deviceID uuid.UUID, code string) (*Session, error)
}

// Common Service Errors
var (
	ErrRequireEventID    = errors.New("missing required event id")
	ErrRequireDeviceID   = errors.New("missing required device id")
	ErrRequireUnlockCode = errors.New("missing required unlock code")
	ErrRepository        = errors.New("unable to query repository")
	ErrEventNotFound     = errors.New("event not found")
	ErrUnlockNotFound    = errors.New("device / unlock code combination not found")
)

// Session holds session details
type Session struct {
	EventCaption  string `json:"event_caption,omitempty"`
	DeviceCaption string `json:"device_caption,omitempty"`
	Token         string `json:"token,omitempty"`
}
