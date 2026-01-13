package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schemaVersion = 1
	defaultDBFile = "app.db"
)

var schema = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS schema_version (
	version INTEGER NOT NULL
);

INSERT INTO schema_version (version) VALUES (%d);
`, schemaVersion)

type DB interface {
	Close() (err error)
}

type SqliteDB struct {
	handle *sqlx.DB
}

func initFrom(from string) (db *SqliteDB, err error) {
	handle, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}

	_, err = _db.handle.Exec(schema)
	if err != nil {
		return
	}

	_db = &SqliteDB{handle: handle}
}

func Init() (db DB, err error) {
	_db, err := initFrom(defaultDBFile)
	if err != nil {
		return
	}

	currentVersion, err := _db.getDbVersion(schema)
	if err != nil {
		return
	}

	if currentVersion > schemaVersion {
		err = fmt.Errorf("forward compatibility not supporte, database version %d is newer than compiled version %d", currentVersion, schemaVersion)
		return
	}

	if currentVersion < schemaVersion {
		err = _db.migrate(currentVersion, schemaVersion)
		if err != nil {
			return
		}
	}

	db = _db

	return
}

func InitFromMemory() (db DB, err error) {
	db, err = initFrom(":memory:")
	return
}

func (db *SqliteDB) getDbVersion() (currentVersion int, err error) {
	err = db.handle.Get(&currentVersion, `SELECT version FROM schema_version LIMIT 1`)
	if err != nil {
		return
	}

	return
}

// Simple ad-hoc migration logic
// production code might be better architecture using https://github.com/golang-migrate/migrate
func (db *SqliteDB) migrate(from int, to int) (err error) {
	if from != to {
		err = migrate(from, to - 1)
		if err != nil {
			return
		}
	} else {
		// down to the current version, nothing to do
		return
	}

	currentVersion, err := db.getDbVersion()
	if err != nil {
		return
	}

	// Sanity check -
	if currentVersion != (to - 1) {
		err = fmt.Errorf("only incremental migration supported")
	}

	switch to {
	default:
		panic("unsupported migration target:", to)
	}


	return
}

func (db *SqliteDB) Close() (err error) {
	err = db.handle.Close()
	return
}
