package sqlclient

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/config"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlclientinternal "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/types"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type SQLClientDependencies struct {
	dig.In
	Loggers module.DependenciesMap[loggercontracts.Logger]
	Tracers module.DependenciesMap[tracercontracts.Tracer]
	Metrics module.DependenciesMap[metricscontracts.Metrics]
}

var SQLClientModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.SQLClientOptions, contracts.SQLClientCore, SQLClientDependencies]{
		Name:             "database.sql-client",
		ConfigKey:        "database.sql-clients",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.SQLClientEnumsDecodeHook()},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.SQLClientOptions,
			dependencies SQLClientDependencies,
		) (contracts.SQLClientCore, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			var tracerInstance tracercontracts.Tracer
			if cfg.EnableTracing {
				tracerInstance, err = dependencies.Tracers.Get(cfg.Tracer)
				if err != nil {
					return nil, err
				}
			}

			var metricsInstance metricscontracts.Metrics
			if cfg.EnableMetrics {
				metricsInstance, err = dependencies.Metrics.Get(cfg.Metrics)
				if err != nil {
					return nil, err
				}
			}

			return sqlclientinternal.CreateSQLClient(
				name,
				cfg,
				loggerInstance,
				tracerInstance,
				metricsInstance,
			)
		},
		HookFn: func(ctx context.Context, name string, instance contracts.SQLClientCore) containercontracts.HookResolveResult {
			return containercontracts.HookResolveResult{
				Name: "database.sql-client:" + name,
				Wait: func(ctx context.Context) error {
					if err := instance.Ping(ctx); err != nil {
						instance.Logger().Errorf("sql-clients %q failed to access database: %v", name, err)
						return fmt.Errorf("sql-client %q ping: %w", name, err)
					}

					<-ctx.Done()
					return nil
				},
				Cleanup: func(ctx context.Context) error {
					if err := instance.Close(ctx); err != nil {
						instance.Logger().Errorf("sql-clients %q close failed with error: %v", name, err)
						return fmt.Errorf("sql-client %q close: %w", name, err)
					}

					instance.Logger().Infof("sql-clients: %q shutdown successfully", name)
					return nil
				},
			}
		},
	},
)
