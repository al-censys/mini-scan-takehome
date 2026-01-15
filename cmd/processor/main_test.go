package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/censys/scan-takehome/pkg/scanning"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/censys/scan-takehome/cmd/processor/db"
)

func loadScansFromJSON(path string) (scans []scanning.Scan, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &scans)
	if err != nil {
		return
	}
	return
}
/*
 * Singel processing
 */
type ProcessScan_MockDB struct {
	db.MockDB
	response string // decoded response
}

func (db *ProcessScan_MockDB) Upsert(r db.Record) (err error) {
	ipLen := len(r.Ip)
	if ipLen == 0 {
		err = fmt.Errorf("invalid IP address length: %d", ipLen)
		return
	}

	if r.Port == 0 {
		err = fmt.Errorf("invalid port number: %d", r.Port)
		return
	}

	db.response = r.Description

	return
}

func TestProcessScan(t *testing.T) {
	ctx := context.Background()

	scan := &scanning.Scan{
		Ip: "10.0.0.1",
		Port: 80,
		Timestamp: time.Now().UnixNano(),
		Service: "HTTP",
		DataVersion: scanning.V2,
		Data: map[string]interface{}{
			"response_str": "service response",
		},
	}

	store := &ProcessScan_MockDB{}

	p := Processor{
		store: store,
	}
	err := p.processScan(ctx, scan)
	assert.NoError(t, err, "p.processScan()")

	// ipv6 supports essentially "free".
	scan.Ip = "fe80::6f64:a408:dc5b:cd46"

	err = p.processScan(ctx, scan)
	assert.NoError(t, err, "p.processScan()")

	// ensure V1 data are base64-decoded
	decoded := "hello world"
	encoded := base64.StdEncoding.EncodeToString([]byte(decoded))
	require.NoError(t, err, "base64.StdEncoding.EncodeString()")

	scan.DataVersion = scanning.V1

	scan.Data = map[string]interface{}{
		"response_bytes_utf8": encoded,
	}

	err = p.processScan(ctx, scan)
	require.NoError(t, err, "p.processScan()")

	assert.Equal(t, decoded, store.response, "p.processScan()")
}

/*
 * Bulk imports within pubsub loop
 */
type BulkImport_MockDB struct {
	db.MockDB
	callCount int
}

func (db *BulkImport_MockDB) Upsert(r db.Record) (err error) {
	db.callCount++
	return
}

type BulkImport_MockSubscription struct {
	scans []scanning.Scan
}

func (sub *BulkImport_MockSubscription) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) (err error) {
	for _, scan := range sub.scans {
		var data []byte
		data, err = json.Marshal(scan)
		if err != nil {
			return
		}

		msg := &pubsub.Message{
			Data: data,
		}

		f(ctx, msg)

		// preempt loop
		select {
		case <-ctx.Done():
			return ctx.Err()
        	default:
		}
	}

	return
}

func TestBulkImport(t *testing.T) {
	scans, err := loadScansFromJSON("test/bulk.json")
	require.NoError(t, err, "loadScansFromJSON()")

	store := &BulkImport_MockDB{}

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)

	subscription := &BulkImport_MockSubscription {
		scans: scans,
	}

	p := Processor{
		store: store,
		subscription: subscription,
	}

	p.Run(ctx)

	assert.Equal(t, len(scans), store.callCount, "invalid number of calls")
}
