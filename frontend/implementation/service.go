package implementation

import (
	"context"

	"github.com/satori/go.uuid"

	"github.com/basvanbeek/opencensus-gokit-example/frontend"
)

type service struct{}

func (s *service) UnlockDevice(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) (frontend.Session, error) {
	return frontend.Session{}, nil
}

func (s *service) GenerateQR(ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string) ([]byte, error) {
	return nil, nil
}

func NewService() frontend.Service {
	return &service{}
}
