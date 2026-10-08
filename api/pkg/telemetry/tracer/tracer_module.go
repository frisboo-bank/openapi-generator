package tracer

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	tracerinternal "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types/tracertype"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type TracerModuleDependencies struct {
	dig.In
}

var TracerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.TracerOptions, contracts.Tracer, TracerModuleDependencies]{
		Name:             "telemetry.tracer",
		ConfigKey:        "telemetry.tracer",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.TracerEnumsDecodeHook()},
		ProviderFn: func(name string, cfg *config.TracerOptions, _ environmentenum.Environment, logger loggercontracts.Logger, _ TracerModuleDependencies) (contracts.Tracer, error) {
			return tracerinternal.CreateTracer(name, cfg, logger)
		},
	},
)

func CreateNoopTracer(name string, logger loggercontracts.Logger) (contracts.Tracer, error) {
	return tracerinternal.CreateTracer(
		name,
		&config.TracerOptions{Type: tracertype.TracerTypes.NOOP},
		logger,
	)
}
