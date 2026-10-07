package migration

import (
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	migrationinternal "frisboo-bank/openapi-generator-service/pkg/database/migration/internal"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

func CreateMigration(
	name string,
	cfg *models.MigrationOptions,
	env environmentenum.Environment,
	logger loggerContracts.Logger,
	extra MigrationDependencies,
) (contracts.Migration, error) {
	if !env.IsDevelopment() {
		return nil, fmt.Errorf("migration can only run in development environment")
	}

	sqlClient, ok := extra.SQLClients[cfg.DBClient]
	if !ok {
		return nil, fmt.Errorf("sql client %q not found for migration %q", cfg.DBClient, name)
	}

	dbClient, ok := sqlClient.(sqlclientContracts.WithDBGetter)
	if !ok {
		return nil, fmt.Errorf("sql client %q does not expose its DB connection", name)
	}

	var adapter contracts.MigrationAdapter
	var err error

	switch cfg.Type {
	case migrationtype.MigrationTypes.GOOSE:
		adapter, err = goose.NewGooseAdapter(name, cfg, dbClient.DB(), logger)
	default:
		err = fmt.Errorf("unsupported Migration type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return &migration{
		adapter: adapter,
	}, nil
}

func CreateMigrationForTests(
	name string,
	db *sql.DB,
	migrationDir string,
	env environmentenum.Environment,
) (contracts.Migration, error) {
	return migrationinternal.CreateMigrationForTests(name, db, migrationDir, env)
}
