package kv

import (
	"context"
	"time"
)

type KeyValueStorer interface {
	// Normally returns a new KV storer that appends a prefix line
	// {dom}.{key} to every operation.
	WithDomain(dom string) KeyValueStorer

	Get(ctx context.Context, key string, out any) (bool, error)
	GetTTL(ctx context.Context, key string, ttl time.Duration, out any) (bool, error)

	// Empty string for non-existent values.
	GetString(ctx context.Context, key string) (string, error)
	// Empty string for non-existent values.
	GetStringTTL(ctx context.Context, key string, ttl time.Duration) (string, error)

	Set(ctx context.Context, key string, value any) error
	SetTTL(ctx context.Context, key string, value any, ttl time.Duration) error

	Delete(ctx context.Context, key string) error
	DeleteExists(ctx context.Context, key string) (bool, error)
	DeleteGet(ctx context.Context, key string, out any) (bool, error)

	// Empty string for non-existent values
	DeleteGetString(ctx context.Context, key string) (string, error)
}
