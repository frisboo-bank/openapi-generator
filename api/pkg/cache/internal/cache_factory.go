package cache

import (
	"frisboo-bank/openapi-generator-service/pkg/cache/config"
	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	"frisboo-bank/openapi-generator-service/pkg/cache/internal/adapters/memory"
	"frisboo-bank/openapi-generator-service/pkg/cache/internal/adapters/redis"
	"frisboo-bank/openapi-generator-service/pkg/cache/types/cachetype"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
)

func CreateCache(
	name string,
	cfg *config.CacheOptions,
	logger loggercontracts.Logger,
) (contracts.Cache, error) {
	var adapter contracts.Cache
	var err error

	switch cfg.Type {
	case cachetype.CacheTypes.REDIS:
		adapter = redis.NewRedisAdapter(name, cfg, logger)
	case cachetype.CacheTypes.MEMORY:
		adapter, err = memory.NewMemoryAdapter(name, cfg, logger)
	default:
		err = syserrors.Newf("no cache-client of type %q exists", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return adapter, nil
}
