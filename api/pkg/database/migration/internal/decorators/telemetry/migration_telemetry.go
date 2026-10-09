package telemetry

import (
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func WrapMigrationForTelemetry(
	name string,
	delegate contracts.Migration,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) contracts.Migration {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("delegate", delegate)

	if tracer == nil && metrics == nil {
		return delegate
	}

	if tracer != nil {
		delegate = decorateMigrationForTracing(name, delegate, tracer)
	}

	if metrics != nil {
		delegate = decorateMigrationForMetrics(name, delegate, metrics)
	}

	return delegate
}
