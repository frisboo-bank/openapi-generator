package migration

import (
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/internal/adapters/goose"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

func CreateMigrationForTests(
	name string,
	db *sql.DB,
	migrationDir string,
	env environmentenum.Environment,
) (contracts.Migration, error) {
	if !env.IsTesting() {
		return nil, fmt.Errorf("migration can only run in testing environment")
	}

	loggerInstance, logErr := logger.CreateNoopLogger("test", environmentenum.Environments.TESTING)
	if logErr != nil {
		return nil, logErr
	}

	adapter, err := goose.NewGooseAdapter(
		name,
		&config.MigrationOptions{
			MigrationsDir: migrationDir,
		},
		db,
		loggerInstance,
	)
	if err != nil {
		return nil, err
	}

	return &migration{adapter}, nil
}

func CreateMigration(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	env environmentenum.Environment,
	logger loggerContracts.Logger,
) (contracts.Migration, error) {
	if !env.IsDevelopment() {
		return nil, fmt.Errorf("migration can only run in development environment")
	}

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

	return &migration{adapter: adapter}, nil
}
