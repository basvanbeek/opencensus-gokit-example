package implementation

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/device"
	"github.com/basvanbeek/opencensus-gokit-example/device/database"
)

// service implements frontend.Service
type service struct {
	repository database.Repository
	logger     log.Logger
}

// NewService creates and returns a new Device service instance
func NewService(rep database.Repository, logger log.Logger) device.Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

// Unlock returns new session data for allowing device to check-in participants.
func (s *service) Unlock(
	ctx context.Context, eventID, deviceID uuid.UUID, unlockCode string,
) (device.Session, error) {
	return device.Session{}, nil
}
