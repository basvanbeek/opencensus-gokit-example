package implementation

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/satori/go.uuid"

	"github.com/basvanbeek/opencensus-gokit-example/frontend"
	"github.com/basvanbeek/opencensus-gokit-example/qr"
)

type service struct {
	logger   log.Logger
	qrClient qr.Service
}

func (s *service) UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (frontend.Session, error) {
	return frontend.Session{}, nil
}

func (s *service) GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error) {
	level.Debug(s.logger).Log("method", "GenerateQR")
	return s.qrClient.Generate(
		ctx, eventID.String()+":"+deviceID.String()+":"+unlockCode, qr.LevelM, 256,
	)
}

// NewService returns a new frontend service
func NewService(qrClient qr.Service, logger log.Logger) frontend.Service {
	return &service{
		qrClient: qrClient,
		logger:   log.With(logger, "client", "QR"),
	}
}
