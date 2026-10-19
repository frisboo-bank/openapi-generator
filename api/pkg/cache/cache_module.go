package cache

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	"frisboo-bank/openapi-generator-service/pkg/cache/models"
	cacheenums "frisboo-bank/openapi-generator-service/pkg/cache/models/enums"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

const CachesModule = "caches"

type CacheDependencies struct {
	dig.In
	Env     environmentenum.Environment
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var CacheModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.CacheOptions, contracts.Cache, CacheDependencies]{
		Name:      "cache",
		ConfigKey: "caches",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			cacheenums.CacheEnumsDecodeHook(),
		},
		ProviderFn: func(
			name string,
			cfg *models.CacheOptions,
			dependencies CacheDependencies,
		) (contracts.Cache, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return CreateCache(name, cfg, loggerInstance, dependencies.Env)
		},
		HookFn: func(name string, instance contracts.Cache) containerContracts.HookResolveResult {
			return containerContracts.HookResolveResult{
				Name: "cache:" + name,
				Start: func(ctx context.Context) error {
					go func() {
						if err := instance.Ping(ctx); err != nil {
							instance.Logger().Fatalf("cache-client %q failed to access backend with error: %v", instance.Name(), err)
						}
					}()

					<-ctx.Done()

					return nil
				},
				Stop: func(ctx context.Context) error {
					if err := instance.Close(); err != nil {
						instance.Logger().Errorf("cache-client %q close failed with error: %v", name, err)
						return err
					}

					instance.Logger().Infof("cache-client: %q shutdown successfully", name)
					return nil
				},
			}
		},
	},
)
