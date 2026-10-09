package metrics

import (
	"context"
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
)

func CreateMetrics(
	name string,
	cfg *config.MetricsOptions,
	ctx context.Context,
	logger loggercontracts.Logger,
) (contracts.Metrics, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var adapter contracts.Metrics
	var err error

	switch cfg.Type {
	case metricstype.MetricsTypes.OPEN_TELEMETRY:
		adapter, err = otel.NewOtelMetricsAdapter(
			name,
			cfg,
			ctx,
			shared.Resource("openapi-generator-service"),
			logger,
		)
	default:
		err = fmt.Errorf("unsupported Metrics type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	return adapter, nil
}
