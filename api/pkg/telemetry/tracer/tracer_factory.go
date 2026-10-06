package tracer

import (
	"context"
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/noop"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
)

func CreateTracer(
	name string,
	cfg *models.TracerOptions,
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
		adapter, err = otel.NewOtelTracerAdapter(name, cfg, ctx, shared.Resource("openapi-generator-service"), logger)
	case tracertype.TracerTypes.NOOP:
		adapter, err = noop.NewNoopTracerAdapter(name, ctx, logger)
	default:
		err = fmt.Errorf("unsupported Tracer type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return &tracer{adapter: adapter}, nil
}

func CreateNoopTracer(
	name string,
	ctx context.Context,
	logger loggercontracts.Logger,
) (contracts.Tracer, error) {
	return CreateTracer(
		name,
		&models.TracerOptions{Type: tracertype.TracerTypes.NOOP},
		ctx,
		logger,
	)
}
