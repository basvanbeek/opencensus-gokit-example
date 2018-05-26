package sqlite

import (
	// external
	"github.com/jmoiron/sqlx"
)

func v1(tx *sqlx.Tx) (err error) {
	// add event table
	if _, err = tx.Exec(`
		CREATE TABLE event (
			id BLOB NOT NULL, tenant_id BLOB NOT NULL, name TEXT NOT NULL,
			PRIMARY KEY(id)
		) WITHOUT ROWID;`,
	); err != nil {
		return
	}

	if _, err = tx.Exec(
		`CREATE UNIQUE INDEX uidx_event_name ON event (tenant_id, lower(name));`,
	); err != nil {
		return
	}

	return
}
