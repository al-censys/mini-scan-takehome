package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/netip"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/censys/scan-takehome/pkg/scanning"

	"github.com/censys/scan-takehome/cmd/processor/db"
)

type Processor struct {
	store db.DB
	storeLock sync.Mutex
	subscription Subscription
}

func (p *Processor) processScan(ctx context.Context, scan *scanning.Scan) (err error) {
	slog.Info("processing", "scan", scan)

	ip, err := netip.ParseAddr(scan.Ip)
	if err != nil {
		return
	}

	// support both ipv4 and ipv6
	ip_data := []byte{}
	if ip.Is4() {
		a := ip.As4()
		ip_data = append(ip_data, a[:]...)
	} else if ip.Is6() {
		a := ip.As16()
		ip_data = append(ip_data, a[:]...)
	} else {
		err = fmt.Errorf("unhandled IP format")
		return
	}

	// while we are at it, we should probably do some sanity checks, but
	// current trivial implemenation lacks error handling infrastructure.
	port := uint16(scan.Port)

	timestamp := time.Unix(scan.Timestamp, 0)

	service := scan.Service

	// TODO - rename to match the "response" name used in the scan
	description := ""
	m := scan.Data.(map[string]interface{})
	switch scan.DataVersion {
	case scanning.V1:
		base64Encoded := m["response_bytes_utf8"].(string)

		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(base64Encoded)
		if err != nil {
			return
		}

		description = string(decoded)
	case scanning.V2:
		description = m["response_str"].(string)
	default:
		err = fmt.Errorf("unknown data version: %d", scan.DataVersion)
		return
	}

	record := db.Record{
		Ip: ip_data,
		Port: port,
		Service: service,
		Timestamp: timestamp,
		Description: description,
	}

	p.storeLock.Lock()
	defer p.storeLock.Unlock()

	err = p.store.Upsert(record)
	if err != nil {
		return
	}

	return
}

func (p *Processor) processMessage(ctx context.Context, m *pubsub.Message) {
	data := m.Data

	var scan scanning.Scan
	err := json.Unmarshal(data, &scan)
	if err != nil {
		slog.Error(err.Error())
		m.Nack()
		return
	}

	err = p.processScan(ctx, &scan)
	if err != nil {
		slog.Error(err.Error())
		m.Nack()
		return
	}

	// acknowledge message
	// at-least-once model - check
	m.Ack()
}

func (p *Processor) Run(ctx context.Context) (err error) {
	fmt.Println(p.subscription)
	err = p.subscription.Receive(ctx, p.processMessage)
	if err != nil {
		return
	}

	return
}


func main() {
	slog.Info("Starting processor...")

	// from cmd/scanner, would have been nice to have a common helper to avoid
	// having to copy/paste
	projectId := flag.String("project", "test-project", "GCP Project ID")
	topicId := flag.String("topic", "scan-topic", "GCP PubSub Topic ID")

	subscriptionId := flag.String("subscription", "scan-subscription", "GCP PubSub Subscription ID")

	ctx := context.Background()

	slog.Info("Connecting to pub/sub...")
	client, err := pubsub.NewClient(ctx, *projectId)
	if err != nil {
		panic(err)
	}

	topic, err := scanning.CreateTopic(ctx, client, *topicId)
	if err != nil {
		panic(err)
	}

	sub := client.Subscription(*subscriptionId)
	exists, err := sub.Exists(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		sub, err = client.CreateSubscription(ctx, *subscriptionId, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
		})
		if err != nil {
			panic(err)
		}
	}

	// Init DB
	store := db.MustInit()
	defer store.Close()

	// let's go...
	p := &Processor{
		store: store,
		subscription: sub,
	}

	err = p.Run(ctx)
	if err != nil {
		panic(err)
	}

	return
}
