package metrics

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models"

	"go.uber.org/dig"
)

type MetricsModuleDependencies struct {
	dig.In
}

var MetricsModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.MetricsOptions, contracts.Metrics, MetricsModuleDependencies]{
		Name:      "metrics",
		ConfigKey: "metrics",
		ProviderFn: func(name string, cfg *models.MetricsOptions, _ environmentEnum.Environment, logger loggercontracts.Logger, _ MetricsModuleDependencies) (contracts.Metrics, error) {
			return CreateMetrics(name, cfg, logger)
		},
	},
)
