package qr

import (
	"context"
)

// RecoveryLevel : Error detection/recovery capacity.
//
// There are several levels of error detection/recovery capacity. Higher levels
// of error recovery are able to correct more errors, with the trade-off of
// increased symbol size.
type RecoveryLevel int

// RecoveryLevel enum
const (
	LevelL RecoveryLevel = iota // Level L: 7% error recovery.
	LevelM                      // Level M: 15% error recovery. Good default choice.
	LevelQ                      // Level Q: 25% error recovery.
	LevelH                      // Level H: 30% error recovery.
)

type Service interface {
	Generate(ctx context.Context, url string, level RecoveryLevel, size int) ([]byte, error)
}
