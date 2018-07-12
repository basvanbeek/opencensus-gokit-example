package implementation

import (
	// stdlib
	"context"
	"database/sql"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kevinburke/go.uuid"
	"golang.org/x/crypto/bcrypt"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/device"
	"github.com/basvanbeek/opencensus-gokit-example/services/device/database"
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
) (*device.Session, error) {
	logger := log.With(s.logger, "method", "Unlock")

	details, err := s.repository.GetDevice(ctx, eventID, deviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			level.Error(logger).Log("err", err)
			return nil, device.ErrRepository
		}
		details = &database.Session{}
	}

	if err = bcrypt.CompareHashAndPassword(
		details.UnlockHash, []byte(unlockCode),
	); err != nil {
		level.Error(logger).Log("err", err)
		return nil, device.ErrUnlockNotFound
	}

	return &device.Session{
		EventCaption:  details.EventCaption,
		DeviceCaption: details.DeviceCaption,
	}, nil
}
