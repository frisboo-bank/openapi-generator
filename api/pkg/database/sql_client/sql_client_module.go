package sqlclient

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	sqlclientenums "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums"
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
	module.MultiInstancesModuleOptions[*models.SQLClientOptions, contracts.SQLClientCore, SQLClientDependencies]{
		Name:      "database.sql-client",
		ConfigKey: "database.sql-clients",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			sqlclientenums.SQLClientEnumsDecodeHook(),
		},
		ProviderFn: func(
			name string,
			cfg *models.SQLClientOptions,
			env environmentenum.Environment,
			logger loggercontracts.Logger,
			extra SQLClientDependencies,
		) (contracts.SQLClientCore, error) {
			ctx := context.Background()

			var err error

			var tracerInstance tracercontracts.Tracer
			if !cfg.EnableTracing {
				tracerInstance, err = tracer.CreateNoopTracer(name, context.Background(), logger)
				if err != nil {
					return nil, err
				}
			} else {
				var ok bool
				tracerInstance, ok = extra.Tracers[cfg.Tracer]
				if !ok {
					return nil, fmt.Errorf("tracer %q not found for sql-client %q", cfg.Tracer, name)
				}
			}

			var metricsInstance metricscontracts.Metrics
			if !cfg.EnableMetrics {
				metricsInstance, err = metrics.CreateNoopMetrics(name, ctx, logger)
				if err != nil {
					return nil, err
				}
			} else {
				var ok bool
				metricsInstance, ok = extra.Metrics[cfg.Metrics]
				if !ok {
					return nil, fmt.Errorf("metrics %q not found for sql-client %q", cfg.Metrics, name)
				}
			}

			return CreateSQLClient(name, cfg, ctx, logger, tracerInstance, metricsInstance)
		},
		HookFn: func(name string, instance contracts.SQLClientCore) containercontracts.HookResolveResult {
			return containercontracts.HookResolveResult{
				Name: "database.sql-client:" + name,
				Wait: func(ctx context.Context) error {
					if err := instance.Ping(); err != nil {
						instance.Logger().Errorf("sql-clients %q failed to access database: %v", name, err)
						return fmt.Errorf("sql-client %q ping: %w", name, err)
					}

					<-ctx.Done()
					return nil
				},
				Cleanup: func(ctx context.Context) error {
					if err := instance.Close(); err != nil {
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
