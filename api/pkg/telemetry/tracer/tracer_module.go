package tracer

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	tracerinternal "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type TracerModuleDependencies struct {
	dig.In
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var TracerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.TracerOptions, contracts.Tracer, TracerModuleDependencies]{
		Name:             "telemetry.tracer",
		ConfigKey:        "telemetry.tracer",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.TracerEnumsDecodeHook()},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.TracerOptions,
			dependencies TracerModuleDependencies,
		) (contracts.Tracer, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return tracerinternal.CreateTracer(ctx, name, cfg, loggerInstance)
		},
	},
)
