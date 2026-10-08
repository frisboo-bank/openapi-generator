package tracer

import (
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/noop"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types/tracertype"
)

func CreateTracer(
	name string,
	cfg *config.TracerOptions,
	logger loggercontracts.Logger,
) (contracts.Tracer, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var adapter contracts.TracerAdapter
	var err error

	switch cfg.Type {
	case tracertype.TracerTypes.OPEN_TELEMETRY:
		adapter, err = otel.NewOtelTracerAdapter(name, cfg, shared.Resource("openapi-generator-service"), logger)
	case tracertype.TracerTypes.NOOP:
		adapter, err = noop.NewNoopTracerAdapter(name, logger)
	default:
		err = fmt.Errorf("unsupported Tracer type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return &tracer{adapter: adapter}, nil
}
