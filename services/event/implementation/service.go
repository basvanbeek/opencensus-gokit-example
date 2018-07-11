package implementation

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event"
	"github.com/basvanbeek/opencensus-gokit-example/services/event/database"
)

// service implements event.Service
type service struct {
	i          int32
	repository database.Repository
	logger     log.Logger
}

// NewService creates and returns a new Event service instance
func NewService(rep database.Repository, logger log.Logger) event.Service {
	return &service{
		i:          0,
		repository: rep,
		logger:     logger,
	}
}

func (s *service) Create(
	ctx context.Context, tenantID uuid.UUID, e event.Event,
) (*uuid.UUID, error) {
	logger := log.With(s.logger, "method", "Create")

	dbEvent := database.Event{
		TenantID: tenantID,
		Name:     e.Name,
	}
	id, err := s.repository.Create(ctx, dbEvent)
	switch err {
	case nil:
		return id, nil
	case database.ErrRepository:
		level.Error(logger).Log("err", err)
		return nil, event.ErrService
	case database.ErrNameExists:
		level.Debug(logger).Log("err", err)
		return nil, event.ErrEventExists
	default:
		level.Error(logger).Log("err", err)
		return nil, event.ErrService
	}
}

func (s *service) Get(
	ctx context.Context, tenantID, id uuid.UUID,
) (*event.Event, error) {
	logger := log.With(s.logger, "method", "Get")

	dbEvent, err := s.repository.Get(ctx, id)
	switch err {
	case nil:
		if !uuid.Equal(dbEvent.TenantID, tenantID) {
			// let's not leak event id's from other tenants.
			return nil, event.ErrNotFound
		}
		return &event.Event{ID: dbEvent.ID, Name: dbEvent.Name}, nil
	case database.ErrRepository:
		level.Error(logger).Log("err", err)
		return nil, event.ErrService
	case database.ErrNotFound:
		level.Debug(logger).Log("err", err)
		return nil, event.ErrNotFound
	default:
		level.Error(logger).Log("err", err)
		return nil, event.ErrService
	}
}

func (s *service) Update(ctx context.Context, tenantID uuid.UUID, e event.Event) error {
	logger := log.With(s.logger, "method", "Update")

	err := s.repository.Update(
		ctx,
		database.Event{
			ID:       e.ID,
			TenantID: tenantID,
			Name:     e.Name,
		},
	)

	switch err {
	case nil:
		return nil
	case database.ErrRepository:
		level.Error(logger).Log("err", err)
		return event.ErrService
	case database.ErrNameExists:
		level.Debug(logger).Log("err", err)
		return event.ErrEventExists
	case database.ErrNotFound:
		level.Debug(logger).Log("err", err)
		return event.ErrNotFound
	default:
		level.Error(logger).Log("err", err)
		return event.ErrService
	}
}

func (s *service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	logger := log.With(s.logger, "method", "Delete")

	if err := s.repository.Delete(ctx, tenantID, id); err != nil {
		level.Error(logger).Log("err", err)
		return event.ErrService
	}
	return nil
}

func (s *service) List(ctx context.Context, tenantID uuid.UUID) ([]*event.Event, error) {
	logger := log.With(s.logger, "method", "List")

	dbEvents, err := s.repository.List(ctx, tenantID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, event.ErrService
	}
	events := make([]*event.Event, 0, len(dbEvents))
	for _, dbEvent := range dbEvents {
		events = append(
			events,
			&event.Event{
				ID:   dbEvent.ID,
				Name: dbEvent.Name,
			},
		)
	}
	return events, nil
}
