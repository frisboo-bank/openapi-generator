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
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type SQLClientDependencies struct {
	dig.In
	Tracers map[string]tracercontracts.Tracer
	Metrics map[string]metricscontracts.Metrics
}

var SQLClientModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.SQLClientOptions, contracts.SQLClientCore, SQLClientDependencies]{
		Name:             "database.sql-client",
		ConfigKey:        "database.sql-clients",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.SQLClientEnumsDecodeHook()},
		ProviderFn: func(
			name string,
			cfg *config.SQLClientOptions,
			env environmentenum.Environment,
			logger loggercontracts.Logger,
			extra SQLClientDependencies,
		) (contracts.SQLClientCore, error) {
			var tracerInstance tracercontracts.Tracer
			var ok bool

			if cfg.EnableTracing {
				tracerInstance, ok = extra.Tracers[cfg.Tracer]
				if !ok {
					return nil, fmt.Errorf("tracer %q not found for migration %q", cfg.Tracer, name)
				}
			}

			var metricsInstance metricscontracts.Metrics
			if cfg.EnableMetrics {
				metricsInstance, ok = extra.Metrics[cfg.Metrics]
				if !ok {
					return nil, fmt.Errorf("metrics %q not found for migration %q", cfg.Metrics, name)
				}
			}

			return sqlclientinternal.CreateSQLClient(name, cfg, logger, tracerInstance, metricsInstance)
		},
		HookFn: func(name string, instance contracts.SQLClientCore) containercontracts.HookResolveResult {
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
