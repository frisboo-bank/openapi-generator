package cache

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/cache/config"
	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	cacheinternal "frisboo-bank/openapi-generator-service/pkg/cache/internal"
	"frisboo-bank/openapi-generator-service/pkg/cache/types"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type CacheDependencies struct {
	dig.In
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var CacheModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.CacheOptions, contracts.Cache, CacheDependencies]{
		Name:             "cache",
		ConfigKey:        "caches",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.CacheEnumsDecodeHook()},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.CacheOptions,
			dependencies CacheDependencies,
		) (contracts.Cache, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return cacheinternal.CreateCache(name, cfg, loggerInstance)
		},
		HookFn: func(ctx context.Context, name string, instance contracts.Cache) containerContracts.HookResolveResult {
			return containerContracts.HookResolveResult{
				Name: "cache:" + name,
				Wait: func(ctx context.Context) error {
					go func() {
						if err := instance.Ping(ctx); err != nil {
							instance.Logger().Fatalf("cache-client %q failed to access backend with error: %v", instance.Name(), err)
						}
					}()

					<-ctx.Done()

					return nil
				},
				Cleanup: func(ctx context.Context) error {
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
