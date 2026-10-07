package metrics

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/config"
	metricsinternal "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal"
	metricsenums "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MetricsModuleDependencies struct {
	dig.In
}

var MetricsModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.MetricsOptions, contracts.Metrics, MetricsModuleDependencies]{
		Name:      "telemetry.metrics",
		ConfigKey: "telemetry.metrics",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			metricsenums.MetricsEnumsDecodeHook(),
		},
		ProviderFn: func(name string, cfg *config.MetricsOptions, _ environmentenum.Environment, logger loggercontracts.Logger, _ MetricsModuleDependencies) (contracts.Metrics, error) {
			return metricsinternal.CreateMetrics(name, cfg, logger)
		},
	},
)
