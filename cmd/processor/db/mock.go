package db

type MockDB struct {
}

// compile-time assertion
var _ DB = (*MockDB)(nil)

func (db *MockDB) Upsert(r Record) (err error) {
	err = IAmNotImplemented()
	return
}

func (db *MockDB) Close() (err error) {
	err = IAmNotImplemented()
	return
}
