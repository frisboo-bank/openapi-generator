package migration

import (
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/internal/adapters/goose"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/internal/decorators/telemetry"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
)

func CreateMigration(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	env environmentenum.Environment,
	logger loggercontracts.Logger,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) (contracts.Migration, error) {
	if !env.IsDevelopment() {
		return nil, fmt.Errorf("migration can only run in development environment")
	}

	return createMigration(name, db, cfg, logger, tracer, metrics)
}

func CreateMigrationForTesting(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	logger loggercontracts.Logger,
) (contracts.Migration, error) {
	return createMigration(name, db, cfg, logger, nil, nil)
}

func createMigration(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	logger loggercontracts.Logger,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) (contracts.Migration, error) {
	var adapter contracts.Migration
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

	if tracer != nil || metrics != nil {
		adapter = telemetry.WrapMigrationForTelemetry(name, adapter, tracer, metrics)
	}

	return adapter, nil
}
