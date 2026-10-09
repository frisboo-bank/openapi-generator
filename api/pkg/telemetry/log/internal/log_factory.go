package log

import (
	"context"

	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/types/logtype"
)

func CreateMetrics(
	ctx context.Context,
	name string,
	cfg *config.LogOptions,
	logger loggercontracts.Logger,
) (contracts.Log, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var adapter contracts.Log
	var err error

	switch cfg.Type {
	case logtype.LogTypes.OPEN_TELEMETRY:
		adapter, err = otel.NewOtelLogAdapter(name, cfg, logger)
	default:
		err = fmt.Errorf("unsupported log type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	return adapter, nil
}
