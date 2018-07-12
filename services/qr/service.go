package qr

import (
	// stdlib
	"context"
	"errors"
)

// ServiceName of this service.
const ServiceName = "qr"

// Service describes our QR service.
type Service interface {
	Generate(ctx context.Context, url string, level RecoveryLevel, size int) ([]byte, error)
}

// QR Service Error descriptions
const (
	ErrorInvalidRecoveryLevel = "invalid recovery level requested"
	ErrorInvalidSize          = "invalid size requested"
	ErrorNoContent            = "content can't be empty"
	ErrorContentTooLarge      = "content size too large"
	ErrorGenerate             = "unable to generate QR"
)

// QR Service Errors
var (
	ErrInvalidRecoveryLevel = errors.New(ErrorInvalidRecoveryLevel)
	ErrInvalidSize          = errors.New(ErrorInvalidSize)
	ErrNoContent            = errors.New(ErrorNoContent)
	ErrContentTooLarge      = errors.New(ErrorContentTooLarge)
	ErrGenerate             = errors.New(ErrorGenerate)
)

// RecoveryLevel : Error detection/recovery capacity.
// See: http://www.qrcode.com/en/about/error_correction.html
type RecoveryLevel int

// RecoveryLevel enum identifying QR Code Error Correction Capability
const (
	LevelL RecoveryLevel = iota // Level L: 7% error recovery.
	LevelM                      // Level M: 15% error recovery.
	LevelQ                      // Level Q: 25% error recovery.
	LevelH                      // Level H: 30% error recovery.
)
