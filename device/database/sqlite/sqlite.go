package sqlite

import (
	// stdlib
	"context"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/device/database"
)

type sqlite struct {
	db     *sqlx.DB
	logger log.Logger
}

// New returns a new Repository backed by SQLite
func New(db *sqlx.DB, logger log.Logger) (database.Repository, error) {
	// run our embedded database versioning logic
	if err := runVersioner(db, logger); err != nil {
		return nil, err
	}

	// return our repository
	return &sqlite{
		db:     db,
		logger: log.With(logger, "rep", "sqlite"),
	}, nil
}

// GetDevice retrieves device information
func (s *sqlite) GetDevice(ctx context.Context, eventID, deviceID uuid.UUID) (*database.Session, error) {
	var session = &database.Session{}

	if err := s.db.QueryRow(`
    SELECT
      e.name as event_caption, d.name as device_caption, d.hash
    FROM event e INNER JOIN device d ON e.id = d.event_id
    WHERE event_id = ?1 AND device_id = ?2;
  `,
		eventID.Bytes(), deviceID.Bytes(),
	).Scan(
		session.EventCaption, session.DeviceCaption, session.UnlockHash,
	); err != nil {
		level.Error(s.logger).Log("err", err.Error())
		return nil, err
	}

	return session, nil
}

// Close implements io.Closer
func (s *sqlite) Close() error {
	return s.db.Close()
}
