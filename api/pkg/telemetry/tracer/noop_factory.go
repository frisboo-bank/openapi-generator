package tracer

import (
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal/adapters/noop"
)

func CreateNoopTracer(
	name string,
	logger loggercontracts.Logger,
) (contracts.Tracer, error) {
	return &tracer{adapter: noop.NewNoopTracerAdapter(name, logger)}, nil
}
