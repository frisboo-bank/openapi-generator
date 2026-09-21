package tracer

import (
	"fmt"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/shared"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/otel"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
)

func CreateTracer(name string, cfg *models.TracerOptions, logger loggercontracts.Logger) (contracts.Tracer, error) {
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	switch cfg.Type {
	case tracertype.TracerTypes.OPEN_TELEMETRY:
		adapter, err := otel.NewOtelTracerAdapter(name, cfg, shared.Resource("openapi-generator-service"), logger)
		if err != nil {
			return nil, err
		}
		return &tracer{adapter: adapter}, nil

	default:
		return nil, fmt.Errorf("unsupported Tracer type: %v", cfg.Type)
	}
}
