package sqlite

import (
	// stdlib
	"context"
	"database/sql"

	// external
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/openxact/versioning"
	"github.com/satori/go.uuid"

	// project
	"github.com/basvanbeek/opencensus-gokit-example/services/event/database"
)

type sqlite struct {
	db     *sqlx.DB
	logger log.Logger
}

// New returns a new Repository backed by SQLite
func New(db *sqlx.DB, logger log.Logger) (database.Repository, error) {
	// run our embedded database versioning logic
	versioner, err := versioning.New(
		db, "ocgokitexample.event", level.Debug(logger), false,
	)
	if err != nil {
		return nil, err
	}
	versioner.Add(1, v1)
	if _, err = versioner.Run(); err != nil {
		return nil, err
	}

	// return our repository
	return &sqlite{db: db, logger: logger}, nil
}

func (s *sqlite) Create(
	ctx context.Context, event database.Event,
) (id *uuid.UUID, err error) {
	// check if we need to create a new UUID
	if uuid.Equal(event.ID, uuid.Nil) {
		event.ID, err = uuid.NewV4()
		if err != nil {
			level.Error(s.logger).Log("err", err)
			return nil, database.ErrRepository
		}
	}

	if _, err = s.db.ExecContext(
		ctx,
		`INSERT INTO event (id, tenant_id, name) VALUES (?, ?, ?)`,
		event.ID.Bytes(), event.TenantID.Bytes(), event.Name,
	); err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			switch sqlErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				level.Debug(s.logger).Log("err", err)
				return nil, database.ErrNameExists
			case sqlite3.ErrConstraintPrimaryKey:
				level.Debug(s.logger).Log("err", err)
				return nil, database.ErrIDExists
			}
		}
		level.Error(s.logger).Log("err", err)
		return nil, database.ErrRepository
	}

	return &event.ID, nil
}

func (s *sqlite) Get(ctx context.Context, id uuid.UUID) (*database.Event, error) {
	event := database.Event{ID: id}

	if err := s.db.QueryRowContext(
		ctx,
		`SELECT tenant_id, name FROM event WHERE id = ?`, id.Bytes(),
	).Scan(&event.TenantID, &event.Name); err != nil {
		if err == sql.ErrNoRows {
			level.Debug(s.logger).Log("err", err)
			return nil, database.ErrNotFound
		}
		level.Error(s.logger).Log("err", err)
		return nil, database.ErrRepository
	}

	return &event, nil
}

func (s *sqlite) Update(ctx context.Context, event database.Event) (err error) {
	var (
		res sql.Result
		cnt int64
	)

	res, err = s.db.ExecContext(
		ctx,
		`UPDATE event SET name = ? WHERE tenant_id = ? AND id = ?`,
		event.Name, event.TenantID.Bytes(), event.ID.Bytes(),
	)
	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				level.Debug(s.logger).Log("err", err)
				return database.ErrNameExists
			}
		}
		level.Error(s.logger).Log("err", err)
		return database.ErrRepository
	}

	cnt, err = res.RowsAffected()
	if err != nil {
		level.Error(s.logger).Log("err", err)
		return database.ErrRepository
	}

	if cnt == 0 {
		level.Debug(s.logger).Log("err", err)
		return database.ErrNotFound
	}

	return
}

func (s *sqlite) Delete(
	ctx context.Context, tenantID uuid.UUID, id uuid.UUID,
) (err error) {
	if _, err = s.db.ExecContext(
		ctx,
		`DELETE FROM event WHERE tenant_id = ? AND id = ?`,
		tenantID.Bytes(), id.Bytes(),
	); err != nil {
		level.Error(s.logger).Log("err", err)
		return database.ErrRepository
	}

	return
}

func (s *sqlite) List(
	ctx context.Context, tenantID uuid.UUID,
) (events []*database.Event, err error) {
	var rows *sql.Rows

	if uuid.Equal(tenantID, uuid.Nil) {
		// listing all events
		rows, err = s.db.QueryContext(
			ctx,
			`SELECT id, tenant_id, name FROM event ORDER BY tenant_id, name`,
		)
	} else {
		// listing owned events
		rows, err = s.db.QueryContext(
			ctx,
			`SELECT id, tenant_id, name FROM event WHERE tenant_id = ? ORDER BY name`,
			tenantID.Bytes(),
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*database.Event, 0), nil
		}
		level.Error(s.logger).Log("err", err)
		return nil, database.ErrRepository
	}

	for rows.Next() {
		var event database.Event
		if err = rows.Scan(
			&event.ID, &event.TenantID, &event.Name,
		); err != nil {
			level.Error(s.logger).Log("err", err)
			return nil, database.ErrRepository
		}
		events = append(events, &event)
	}

	return events, nil
}
