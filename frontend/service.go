package frontend

import (
	"context"
	"errors"

	"github.com/satori/go.uuid"
)

// NotFoundError is used for various not found errors
type NotFoundError error

// Common Service Errors
var (
	ErrEventNotFound  NotFoundError = errors.New("event not found")
	ErrUnlockNotFound NotFoundError = errors.New("device / unlock code combination not found")
)

type Service interface {
	UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (Session, error)
	GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error)
}

type Session struct {
	EventID       uuid.UUID
	DeviceID      uuid.UUID
	EventCaption  string
	DeviceCaption string
}
