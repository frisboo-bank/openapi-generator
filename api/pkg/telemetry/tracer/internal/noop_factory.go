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
	adapter, err := noop.NewNoopTracerAdapter(name, logger)
	if err != nil {
		return nil, err
	}
	return &tracer{adapter: adapter}, nil
}
