package cache

import (
	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	"frisboo-bank/openapi-generator-service/pkg/cache/internal/adapters/memory"
	"frisboo-bank/openapi-generator-service/pkg/cache/internal/adapters/redis"
	"frisboo-bank/openapi-generator-service/pkg/cache/models"
	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
)

func CreateCache(
	name string,
	cfg *models.CacheOptions,
	logger loggercontracts.Logger,
	env environmentEnum.Environment,
) (contracts.Cache, error) {
	var adapter contracts.CacheAdapter
	var err error

	switch cfg.Type {
	case cachetype.CacheTypes.REDIS:
		adapter = redis.NewRedisAdapter(cfg, logger, env)
	case cachetype.CacheTypes.MEMORY:
		adapter, err = memory.NewMemoryAdapter(cfg, logger, env)
	default:
		err = syserrors.Newf("no cache-client of type %q exists", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return &cacheFacade{
		name:    name,
		adapter: adapter,
	}, nil
}
