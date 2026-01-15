package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultDB = "data/processor.db"
	latestVersion = 2

	// Initial DB schema, updated on db creation through migrations
	schema = `
CREATE TABLE IF NOT EXISTS schema_version (
	version INTEGER NOT NULL
);

INSERT INTO schema_version (version) VALUES (1);
`
)

type SqliteDB struct {
	handle *sqlx.DB
}

/*
 * Boilerplate init logic
 */
func openFrom(from string) (db *SqliteDB, err error) {
	handle, err := sqlx.Open("sqlite3", from)
	if err != nil {
		return
	}

	db = &SqliteDB{handle: handle}

	return
}

func initFrom(from string) (db *SqliteDB, err error) {
	_db, err := openFrom(from)
	if err != nil {
		return
	}

	_, err = _db.handle.Exec(schema)
	if err != nil {
		return
	}

	db = _db

	return
}

func Init() (db DB, err error) {
	exists, err := fileExists(defaultDB)
	if err != nil {
		return
	}

	var _db *SqliteDB
	if exists {
		_db, err = openFrom(defaultDB)
	} else {
		_db, err = initFrom(defaultDB)
		if err != nil {
			return
		}
	}

	db, err = _db.init()

	return
}

func MustInit() (db DB) {
	db, err := Init()
	if err != nil {
		panic(err)
	}

	return
}

// for tests, no need to export
func initFromMemory() (db DB, err error) {
	_db, err := initFrom(":memory:")
	if err != nil {
		return
	}

	db, err = _db.init()

	return
}

func (_db *SqliteDB) init() (db DB, err error) {
	// from now on, any schema upgrade is done through migrations
	currentVersion, err := _db.getDbVersion()
	if err != nil {
		return
	}

	if currentVersion > latestVersion {
		err = fmt.Errorf("forward compatibility not supported, database version %d is newer than compiled version %d", currentVersion, latestVersion)
		return
	}

	if currentVersion < latestVersion {
		err = _db.migrate(currentVersion, latestVersion)
		if err != nil {
			return
		}
	}

	db = _db

	return
}

/*
 * Versioning and migration logic
 */
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
		err = db.migrate(from, to - 1)
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
	case 2:
		err = db.migrateToV2()
		if err != nil {
			return
		}
	default:
		panic("unsupported migration target")
	}


	return
}
/*
 * Migration summary:
 * - add the main service records table
 */
func (db *SqliteDB) migrateToV2() (err error) {
	update := `
CREATE TABLE services (
    id          INTEGER PRIMARY KEY,
    ip          BLOB    NOT NULL,
    port        INTEGER NOT NULL,
    service     TEXT    NOT NULL,
    description TEXT,
    timestamp   DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (ip, port, service)
);

CREATE TRIGGER update_services_updated_at
AFTER UPDATE ON services
FOR EACH ROW
BEGIN
    UPDATE services
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
`
	_, err = db.handle.Exec(update)
	if err != nil {
		return
	}

	err = db.migrationIncrementVersion()
	if err != nil {
		return
	}

	return
}

func (db *SqliteDB) migrationIncrementVersion() (err error) {
	_, err = db.handle.Exec(`
UPDATE schema_version
SET version = version + 1
`)

	return
}

func (db *SqliteDB) Upsert(r Record) (err error) {
	/* Requirements were to keep a single record on conflicting triplets,
	 * limit this implementation as such.
	 *
	 * It might be cheaper to run a blind INSERT and either have a different
	 * service / worker clean up afterward, or limiting to the latest
	 * `timestamp' on SELECT, depending on runtime behavior. There might
	 * also be a compliance requirements to keep all the records and merely
	 * mark them outdated, but either case would require more insight and
	 * preemptive discussions on the exact required behavior.
	 */
	query := `

INSERT INTO services (ip, port, service, description, timestamp)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(ip, port, service)
DO UPDATE SET
    description = excluded.description,
    timestamp  = excluded.timestamp
WHERE excluded.timestamp > services.timestamp;
`
	// Only precise down to ms. precision, consider switching the
	// `timestamp' field to INTEGER64 and record via .UnitNano() and
	// "should"  be fine until year 2262.
	timestamp := r.Timestamp.Format("2006-01-02T15:04:05.000Z")

	_, err = db.handle.Exec(query, r.Ip, r.Port, r.Service, r.Description, timestamp)
	if err != nil {
		return
	}

	return
}

// for testing...
func (db *SqliteDB) dumpServices() (records Records, err error) {
	query := "SELECT ip, port, service, description, timestamp FROM services;"

	rows, err := db.handle.Queryx(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r Record
		err = rows.StructScan(&r)
		if err != nil {
			return
		}

		records = append(records, r)
	}

	return
}

func (db *SqliteDB) Close() (err error) {
	err = db.handle.Close()
	return
}


