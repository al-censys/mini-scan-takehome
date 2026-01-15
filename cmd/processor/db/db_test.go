package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	db, err := initFromMemory()
	defer db.Close()

	assert.NoError(nil, err, "initFromMemory()")
	assert.NotNil(nil, db, "initFromMemory()")
}

func getRecord() (r Record) {
	r = Record{
		Ip: []byte{1, 2, 3, 4},
		Port: 80,
		Service: "HTTP",
		Description: "dummy description",
		Timestamp: time.UnixMilli(time.Now().UnixMilli()).UTC(),
	}

	return
}

func TestAddEntry(t *testing.T) {
	r := getRecord()

	db, err := initFromMemory()
	require.NoError(t, err, "initFromMemory()")

	err = db.Upsert(r)
	assert.NoError(t, err, "db.Record()")

	// sanity check
	records, err := db.(*SqliteDB).dumpServices()
	require.NoError(t, err, "db.(*SqliteDB).dumpServices()")

	assert.Equal(t, 1, len(records), "wrong number of records")
	assert.True(t, reflect.DeepEqual(r.Timestamp, records[0].Timestamp), "bad record extracted")
}

// XXX same (ip, port, service) triplet, second record is newer, timestamp will be updated
func TestSameTripletNewRecordUpdate(t *testing.T) {
	r0 := getRecord()

	r1 := r0
	r1.Timestamp = r1.Timestamp.Add(24 * time.Hour)

	db, err := initFromMemory()
	require.NoError(t, err, "initFromMemory()")

	err = db.Upsert(r0)
	assert.NoError(t, err, "db.Record()")

	err = db.Upsert(r1)
	assert.NoError(t, err, "db.Record()")

	// sanity check
	records, err := db.(*SqliteDB).dumpServices()
	require.NoError(t, err, "db.(*SqliteDB).dumpServices()")

	assert.Equal(t, 1, len(records), "wrong number of records")
	assert.Equal(t, r1.Timestamp, records[0].Timestamp, "incorrect record timestamp")
}

// XXX same as above, but will not update the record from the past
func TestSameTripletOutOfOrderNewRecordIgnored(t *testing.T) {
	r0 := getRecord()

	r1 := r0
	r1.Timestamp = r1.Timestamp.Add(-24 * time.Hour)

	db, err := initFromMemory()
	require.NoError(t, err, "initFromMemory()")

	err = db.Upsert(r0)
	assert.NoError(t, err, "db.Record()")

	err = db.Upsert(r1)
	assert.NoError(t, err, "db.Record()")

	// sanity check
	records, err := db.(*SqliteDB).dumpServices()
	require.NoError(t, err, "db.(*SqliteDB).dumpServices()")

	assert.Equal(t, 1, len(records), "wrong number of records")
	assert.Equal(t, r0.Timestamp, records[0].Timestamp, "incorrect record timestamp")
}
