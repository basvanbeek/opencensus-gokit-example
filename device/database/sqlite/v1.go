package sqlite

import (
	// external
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

func v1(tx *sqlx.Tx) (err error) {
	// add event table
	if _, err = tx.Exec(`
    CREATE TABLE event (
      id BLOB NOT NULL, name TEXT NOT NULL, PRIMARY KEY(id)
    ) WITHOUT ROWID;
  `); err != nil {
		return
	}

	// add demo events
	var (
		insEvt = `INSERT INTO event(id, name) VALUES(?1, ?2);`
		evt1   = uuid.NewV5(uuid.NamespaceOID, "device.event.1").Bytes()
		evt2   = uuid.NewV5(uuid.NamespaceOID, "device.event.2").Bytes()
	)

	if err = insert(tx, insEvt, evt1, "OpenCensus Spring Marathon"); err != nil {
		return
	}
	if err = insert(tx, insEvt, evt2, "Go kit Turkey Trot"); err != nil {
		return
	}

	// add device table
	if _, err = tx.Exec(`
    CREATE TABLE device (
      id BLOB NOT NULL, event_id BLOB NOT NULL, name TEXT NOT NULL,
      PRIMARY KEY(id)
    ) WITHOUT ROWID;
  `); err != nil {
		return
	}

	// add demo devices
	var (
		insDev = `INSERT INTO device(id, event_id, name) VALUES (?1, ?2, ?3);`
		dev1_1 = uuid.NewV5(uuid.NamespaceOID, "device.device.1.1").Bytes()
		dev1_2 = uuid.NewV5(uuid.NamespaceOID, "device.device.1.2").Bytes()
		dev2_1 = uuid.NewV5(uuid.NamespaceOID, "device.device.2.1").Bytes()
		dev2_2 = uuid.NewV5(uuid.NamespaceOID, "device.device.2.2").Bytes()
	)

	if err = insert(tx, insDev, evt1, dev1_1, "scanner #1"); err != nil {
		return
	}
	if err = insert(tx, insDev, evt1, dev1_2, "scanner #2"); err != nil {
		return
	}
	if err = insert(tx, insDev, evt2, dev2_1, "scanner #1"); err != nil {
		return
	}
	if err = insert(tx, insDev, evt2, dev2_2, "scanner #2"); err != nil {
		return
	}

	return nil
}
