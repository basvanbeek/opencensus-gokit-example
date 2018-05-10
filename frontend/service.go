package frontend

import (
	"context"
	"errors"

	"github.com/satori/go.uuid"
)

// Service describes our Frontend service.
type Service interface {
	UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (Session, error)
	GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error)
}

// Common Service Errors
var (
	ErrEventNotFound  = errors.New("event not found")
	ErrUnlockNotFound = errors.New("device / unlock code combination not found")
)

// Session holds session details
type Session struct {
	EventID       uuid.UUID `json:"event_id,omitempty"`
	EventCaption  string    `json:"event_caption,omitempty"`
	DeviceID      uuid.UUID `json:"device_id,omitempty"`
	DeviceCaption string    `json:"device_caption,omitempty"`
	Token         string    `json:"token,omitempty"`
}
