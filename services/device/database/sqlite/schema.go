package sqlite

import (
	// external
	"github.com/jmoiron/sqlx"
)

func v1(tx *sqlx.Tx) (err error) {
	// add device table
	if _, err = tx.Exec(`
    CREATE TABLE device (
      id BLOB NOT NULL, event_id BLOB NOT NULL, name TEXT NOT NULL,
      hash BLOB NOT NULL, PRIMARY KEY(id)
    ) WITHOUT ROWID;
  `); err != nil {
		return
	}

	return nil
}
