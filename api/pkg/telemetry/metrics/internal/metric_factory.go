package metrics

import (
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal/adapters/noop"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
)

func CreateMetrics(
	name string,
	cfg *config.MetricsOptions,
	logger loggercontracts.Logger,
) (contracts.Metrics, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var adapter contracts.MetricsAdapter
	var err error

	switch cfg.Type {
	case metricstype.MetricsTypes.OPEN_TELEMETRY:
		adapter, err = otel.NewOtelMetricsAdapter(name, cfg, shared.Resource("openapi-generator-service"), logger)
	case metricstype.MetricsTypes.NOOP:
		adapter, err = noop.NewNoopMetricsAdapter(name, logger)
	default:
		err = fmt.Errorf("unsupported Metrics type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	return &metrics{adapter: adapter}, nil
}
