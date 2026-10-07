package metrics

import (
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal/adapters/noop"
)

func CreateNoopMetrics(
	name string,
	logger loggercontracts.Logger,
) (contracts.Metrics, error) {
	return &metrics{adapter: noop.NewNoopMetricsAdapter(name, logger)}, nil
}
