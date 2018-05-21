package sqlite

import (
	// stdlib

	"errors"

	// external
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/openxact/versioning"
)

func runVersioner(db *sqlx.DB, log log.Logger) error {
	versioner, err := versioning.New(db, "ocgokitexample.device", log, false)
	if err != nil {
		return err
	}
	versioner.Add(1, v1)

	_, err = versioner.Run()
	return err
}

func insert(tx *sqlx.Tx, query string, args ...interface{}) error {
	res, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt != 1 {
		return errors.New("expected successful insert")
	}

	return nil
}
