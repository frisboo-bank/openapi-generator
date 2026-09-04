package contracts

import (
	"context"
	"time"

	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type (
	Cache interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string, expiration time.Duration) error
		Del(ctx context.Context, key ...string) error
		Exists(ctx context.Context, key ...string) (int64, error)
		Ping(ctx context.Context) error
		Close() error
		Name() string
		Type() cachetype.CacheType
		Logger() loggerContracts.Logger
	}
)
