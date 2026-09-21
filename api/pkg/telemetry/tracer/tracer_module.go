package tracer

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models"

	"go.uber.org/dig"
)

type TracerModuleDependencies struct {
	dig.In
}

var TracerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.TracerOptions, contracts.Tracer, TracerModuleDependencies]{
		Name:      "tracer",
		ConfigKey: "tracer",
		ProviderFn: func(name string, cfg *models.TracerOptions, _ environmentEnum.Environment, logger loggercontracts.Logger, _ TracerModuleDependencies) (contracts.Tracer, error) {
			return CreateTracer(name, cfg, logger)
		},
	},
)
