package migration

import (
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	migrationtelemetry "frisboo-bank/openapi-generator-service/pkg/database/migration/internal/decorators/telemetry"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/internal/adapters/goose"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
)

func CreateMigration(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	env environmentenum.Environment,
	logger loggerContracts.Logger,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) (contracts.Migration, error) {
	if !env.IsDevelopment() {
		return nil, fmt.Errorf("migration can only run in development environment")
	}

	return createMigration(name, db, cfg, logger, tracer, metrics)
}

// CreateMigrationForTesting is the test-only path: it skips the development
// environment gate so migrations can run under the TESTING environment.
func CreateMigrationForTesting(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	logger loggerContracts.Logger,
) (contracts.Migration, error) {
	return createMigration(name, db, cfg, logger, nil, nil)
}

func createMigration(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	logger loggerContracts.Logger,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) (contracts.Migration, error) {
	var adapter contracts.MigrationAdapter
	var err error

	switch cfg.Type {
	case migrationtype.MigrationTypes.GOOSE:
		adapter, err = goose.NewGooseAdapter(name, cfg, db, logger)
	default:
		err = fmt.Errorf("unsupported Migration type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	// Decorate the adapter with tracing and metrics before the facade wraps
	// it, so the public contracts.Migration surface stays unchanged.
	// Each wrapper is applied only when its instrument is available — no
	// noop fallback is needed; callers that have no tracer or metrics simply
	// skip that wrapper.
	if tracer != nil {
		adapter = migrationtelemetry.WrapMigrationAdapterForTracing(name, adapter, tracer)
	}
	if metrics != nil {
		adapter = migrationtelemetry.WrapMigrationAdapterForMetrics(name, adapter, metrics)
	}

	return &migration{
		adapter: adapter,
	}, nil
}
