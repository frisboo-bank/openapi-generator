package metrics

import (
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	metricsinternal "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal"
)

func CreateNoopMetrics(
	name string,
	logger loggercontracts.Logger,
) (contracts.Metrics, error) {
	return metricsinternal.CreateNoopMetrics(name, logger)
}
