package cache

import (
	"context"
)

type CacheClient interface {
	Ping(ctx context.Context) error
	Options() *ValkeyOptions
	Close() error
}

func NewCacheClient() (CacheClient, error) {
	return NewValkeyCacheClient()
}
