package implementation

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/skip2/go-qrcode"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/qr"
)

// service implements qr.Service
type service struct {
	logger log.Logger
}

// NewService creates and returns a new QR service instance
func NewService(logger log.Logger) qr.Service {
	return &service{
		logger: logger,
	}
}

// Generate returns a new QR code image based on provided details
func (s *service) Generate(
	ctx context.Context, url string, recLevel qr.RecoveryLevel, size int,
) ([]byte, error) {
	var (
		logger = log.With(s.logger, "method", "Generate")
	)
	// test for valid input
	if len(url) == 0 {
		return nil, qr.ErrNoContent
	}
	if recLevel < qr.LevelL || recLevel > qr.LevelH {
		return nil, qr.ErrInvalidRecoveryLevel
	}
	if size > 4096 {
		return nil, qr.ErrInvalidSize
	}
	// do the actual work
	b, err := qrcode.Encode(url, qrcode.RecoveryLevel(recLevel), size)
	if err != nil {
		// actual qrcode lib error... log it...
		level.Error(logger).Log("err", err)

		// consumer of this api gets a generic error returned so we don't leak
		// implementation details upstream
		switch err.Error() {
		case "content too long to encode", "length too long to be represented":
			err = qr.ErrContentTooLarge
		default:
			err = qr.ErrGenerate
		}
	}

	return b, err
}
