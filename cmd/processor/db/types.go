package db

import (
	"time"
)

type DB interface {
	Upsert(r Record) (err error)
	Close() (err error)
}

type Record struct {
	Ip   []byte `db:"ip"`
	Port uint16 `db:"port"`
	Service string `db:"service"`
	Description string `db:"description"`
	Timestamp time.Time `db:"timestamp"`
}

type Records []Record
