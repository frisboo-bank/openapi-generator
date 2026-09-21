package metrics

import (
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models"
	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
)

func CreateMetrics(name string, cfg *models.MetricsOptions, logger loggercontracts.Logger) (contracts.Metrics, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	switch cfg.Type {
	case metricstype.MetricsTypes.OPEN_TELEMETRY:
		adapter, err := otel.NewOtelMetricsAdapter(name, cfg, shared.Resource("openapi-generator-service"), logger)
		if err != nil {
			return nil, err
		}
		return &metrics{adapter: adapter}, nil

	default:
		return nil, fmt.Errorf("unsupported Metrics type: %v", cfg.Type)
	}
}
