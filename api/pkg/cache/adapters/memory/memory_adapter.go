package memory

import (
	"context"
	"fmt"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	"frisboo-bank/openapi-generator-service/pkg/cache/models"
	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/Yiling-J/theine-go"
)

var _ contracts.Cache = (*memoryAdapter)(nil)

type memoryAdapter struct {
	name   string
	cache  *theine.Cache[string, string]
	logger loggerContracts.Logger
}

func NewMemoryAdapter(
	name string,
	cfg *models.CacheOptions,
	logger loggerContracts.Logger,
	env environmentEnum.Environment,
) (contracts.Cache, error) {
	cache, err := theine.NewBuilder[string, string](cfg.MaxEntries).Build()
	if err != nil {
		return nil, err
	}

	return &memoryAdapter{
		name:   name,
		cache:  cache,
		logger: logger,
	}, nil
}

func (m *memoryAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if ok := m.cache.SetWithTTL(key, value, 1, expiration); !ok {
		return fmt.Errorf("cache size exceeded")
	}
	return nil
}

func (m *memoryAdapter) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.cache.Get(key)
	if !ok {
		return "", nil
	}
	return val, nil
}

func (m *memoryAdapter) Del(ctx context.Context, key ...string) error {
	for _, k := range key {
		m.cache.Delete(k)
	}
	return nil
}

func (m *memoryAdapter) Exists(ctx context.Context, key ...string) (int64, error) {
	var count int64
	for _, k := range key {
		if _, ok := m.cache.Get(k); ok {
			count++
		}
	}
	return count, nil
}

func (m *memoryAdapter) Ping(ctx context.Context) error {
	return nil
}

func (m *memoryAdapter) Close() error {
	m.cache.Close()
	return nil
}

func (m *memoryAdapter) Logger() loggerContracts.Logger {
	return m.logger
}

func (m *memoryAdapter) Name() string {
	return m.name
}

func (m *memoryAdapter) Type() cachetype.CacheType {
	return cachetype.CacheTypes.MEMORY
}
