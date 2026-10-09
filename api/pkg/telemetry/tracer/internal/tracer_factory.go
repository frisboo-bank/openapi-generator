package tracer

import (
	"context"
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types/tracertype"
)

func CreateTracer(
	ctx context.Context,
	name string,
	cfg *config.TracerOptions,
	ctx context.Context,
	logger loggercontracts.Logger,
) (contracts.Tracer, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var adapter contracts.Tracer
	var err error

	switch cfg.Type {
	case tracertype.TracerTypes.OPEN_TELEMETRY:
		adapter, err = otel.NewOtelTracerAdapter(
			ctx,
			name,
			cfg,
			shared.Resource("openapi-generator-service"),
			logger,
		)
	default:
		err = fmt.Errorf("unsupported Tracer type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	return adapter, nil
}
