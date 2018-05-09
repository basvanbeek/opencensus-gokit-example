package qr

import (
	"context"
	"errors"
)

// Common Errors for QR Service
var (
	ErrInvalidRecoveryLevel = errors.New("invalid recovery level requested")
	ErrInvalidSize          = errors.New("invalid size requested")
	ErrGenerate             = errors.New("unable to generate QR")
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

type Service interface {
	Generate(ctx context.Context, url string, level RecoveryLevel, size int) ([]byte, error)
}
