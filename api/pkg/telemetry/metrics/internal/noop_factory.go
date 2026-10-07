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
	adapter, err := noop.NewNoopMetricsAdapter(name, logger)
	if err != nil {
		return nil, err
	}
	return &metrics{adapter: adapter}, nil
}
