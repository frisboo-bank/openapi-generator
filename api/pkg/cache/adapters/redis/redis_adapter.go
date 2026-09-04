package redis

import (
	"context"
	"fmt"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	"frisboo-bank/openapi-generator-service/pkg/cache/models"
	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	vendorRedis "github.com/redis/go-redis/v9"
)

var _ contracts.Cache = (*redisAdapter)(nil)

type redisAdapter struct {
	name   string
	cache  *vendorRedis.Client
	logger loggerContracts.Logger
}

func NewRedisAdapter(
	name string,
	cfg *models.CacheOptions,
	logger loggerContracts.Logger,
	env environmentEnum.Environment,
) contracts.Cache {
	srv := vendorRedis.NewClient(&vendorRedis.Options{
		Addr:         cfg.Address(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	return &redisAdapter{
		name:   name,
		cache:  srv,
		logger: logger,
	}
}

func (r *redisAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if err := r.cache.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis set %q: %w", key, err)
	}
	return nil
}

func (r *redisAdapter) Get(ctx context.Context, key string) (string, error) {
	val, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		if err == vendorRedis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis get %q: %w", key, err)
	}
	return val, nil
}

func (r *redisAdapter) Del(ctx context.Context, key ...string) error {
	if err := r.cache.Del(ctx, key...).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}

func (r *redisAdapter) Exists(ctx context.Context, key ...string) (int64, error) {
	return r.cache.Exists(ctx, key...).Result()
}

func (r *redisAdapter) Ping(ctx context.Context) error {
	if err := r.cache.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}

func (r *redisAdapter) Close() error {
	if err := r.cache.Close(); err != nil {
		return fmt.Errorf("redis close: %w", err)
	}
	return nil
}

func (r *redisAdapter) Type() cachetype.CacheType {
	return cachetype.CacheTypes.REDIS
}

func (r *redisAdapter) Logger() loggerContracts.Logger {
	return r.logger
}

func (r *redisAdapter) Name() string {
	return r.name
}
