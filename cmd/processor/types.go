package main

import (
	"context"

	"cloud.google.com/go/pubsub"
)

/*
 * Unfortunately, upstream pubsub does not provide a proper mockable interface.
 * Define our own with associated compile-time assertion to ensure upstream API change won't
 * fail silently.
 *
 * In production code, this should certainly encompass a larger set of the
 * upstream tyes. This is the minimum to get a somewhat testable implementation.
 */
type Subscription interface {
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

var _ Subscription = (*pubsub.Subscription)(nil)

