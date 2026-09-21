package cache

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

var _ contracts.Cache = (*cacheFacade)(nil)

type cacheFacade struct {
	name    string
	adapter contracts.CacheAdapter
}

func (c *cacheFacade) Close() error {
	return c.adapter.Close()
}

func (c *cacheFacade) Del(ctx context.Context, key ...string) error {
	return c.adapter.Del(ctx, key...)
}

func (c *cacheFacade) Exists(ctx context.Context, key ...string) (int64, error) {
	return c.adapter.Exists(ctx, key...)
}

func (c *cacheFacade) Get(ctx context.Context, key string) (string, error) {
	return c.adapter.Get(ctx, key)
}

func (c *cacheFacade) Ping(ctx context.Context) error {
	return c.adapter.Ping(ctx)
}

func (c *cacheFacade) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.adapter.Set(ctx, key, value, expiration)
}

func (c *cacheFacade) Name() string                   { return c.name }
func (c *cacheFacade) Type() cachetype.CacheType      { return c.adapter.Type() }
func (c *cacheFacade) Logger() loggercontracts.Logger { return c.adapter.Logger() }
