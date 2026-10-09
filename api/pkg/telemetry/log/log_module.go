package log

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/types"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	loginternal "frisboo-bank/openapi-generator-service/pkg/telemetry/log/internal"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type LogModuleDependencies struct {
	dig.In
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var LogModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.LogOptions, contracts.Log, LogModuleDependencies]{
		Name:             "telemetry.log",
		ConfigKey:        "telemetry.log",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.LogEnumsDecodeHook()},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.LogOptions,
			dependencies LogModuleDependencies,
		) (contracts.Log, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return loginternal.CreateMetrics(ctx, name, cfg, loggerInstance)
		},
	},
)
