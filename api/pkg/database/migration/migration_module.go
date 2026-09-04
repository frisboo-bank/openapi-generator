package migration

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/adapters/goose"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/models"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	SQLClients map[string]sqlclientContracts.SQLClientCore
}

var MigrationModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.MigrationOptions, contracts.Migration, MigrationDependencies]{
		Name:      "migration",
		ConfigKey: "migration",
		ProviderFn: func(
			name string,
			cfg *models.MigrationOptions,
			env environmentEnum.Environment,
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
				return nil, fmt.Errorf("sql client %q does not support DB retrieval", name)
			}

			switch cfg.Type {
			case migrationtype.MigrationTypes.GOOSE:
				return goose.NewGooseAdapter(name, cfg, dbClient.DB(), logger)
			default:
				return nil, fmt.Errorf("unsupported Migration type: %v", cfg.Type)
			}
		},
	},
)
