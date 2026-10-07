package tracer

import (
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	tracerinternal "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/internal"
)

func CreateNoopTracer(
	name string,
	logger loggercontracts.Logger,
) (contracts.Tracer, error) {
	return tracerinternal.CreateNoopTracer(name, logger)
}
