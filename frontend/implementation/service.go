package implementation

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/qr"
)

// service implements frontend.Service
type service struct {
	qrClient qr.Service
	logger   log.Logger
}

// NewService creates and returns a new Frontend service instance
func NewService(qrClient qr.Service, logger log.Logger) frontend.Service {
	return &service{
		qrClient: qrClient,
		logger:   logger,
	}
}

// Unlockdevice returns a new session for allowing device to check-in participants.
func (s *service) UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (frontend.Session, error) {
	return frontend.Session{}, nil
}

// Generate returns a new QR code device unlock image based on the provided details.
func (s *service) GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error) {
	level.Debug(s.logger).Log("method", "GenerateQR")
	return s.qrClient.Generate(
		ctx, eventID.String()+":"+deviceID.String()+":"+unlockCode, 10, -1,
	)
}
