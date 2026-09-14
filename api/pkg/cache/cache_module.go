package cache

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	"frisboo-bank/openapi-generator-service/pkg/cache/models"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"go.uber.org/dig"
)

const CachesModule = "caches"

type CacheDependencies struct {
	dig.In
}

var CacheModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.CacheOptions, contracts.Cache, CacheDependencies]{
		Name:      "cache",
		ConfigKey: "caches",
		ProviderFn: func(
			name string,
			cfg *models.CacheOptions,
			env environmentEnum.Environment,
			logger loggerContracts.Logger,
			extra CacheDependencies,
		) (contracts.Cache, error) {
			return CreateCache(name, cfg, logger, env)
		},
		HookFn: func(name string, instance contracts.Cache) containerContracts.HookResolveResult {
			return containerContracts.HookResolveResult{
				Name: "cache:" + name,
				Wait: func(ctx context.Context) error {
					go func() {
						if err := instance.Ping(ctx); err != nil {
							instance.Logger().
								Fatalf("cache-client %q failed to access backend with error: %v", instance.Name(), err)
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
