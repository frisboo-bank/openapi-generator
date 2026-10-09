package metrics

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	metricsinternal "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MetricsModuleDependencies struct {
	dig.In
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var MetricsModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.MetricsOptions, contracts.Metrics, MetricsModuleDependencies]{
		Name:             "telemetry.metrics",
		ConfigKey:        "telemetry.metrics",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.MetricsEnumsDecodeHook()},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.MetricsOptions,
			dependencies MetricsModuleDependencies,
		) (contracts.Metrics, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return metricsinternal.CreateMetrics(ctx, name, cfg, loggerInstance)
		},
	},
)
