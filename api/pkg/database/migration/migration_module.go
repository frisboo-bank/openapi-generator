package migration

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/models"
	migrationenums "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	SQLClients map[string]sqlclientcontracts.SQLClientCore
}

var MigrationModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.MigrationOptions, contracts.Migration, MigrationDependencies]{
		Name:      "database.migration",
		ConfigKey: "database.migration",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			migrationenums.MigrationEnumsDecodeHook(),
		},
		ProviderFn: func(
			name string,
			cfg *models.MigrationOptions,
			env environmentenum.Environment,
			logger loggercontracts.Logger,
			extra MigrationDependencies,
		) (contracts.Migration, error) {
			if !env.IsDevelopment() {
				return nil, fmt.Errorf("migration can only run in development environment")
			}

			sqlClient, ok := extra.SQLClients[cfg.DBClient]
			if !ok {
				return nil, fmt.Errorf("sql client %q not found for migration %q", cfg.DBClient, name)
			}

			dbClient, ok := sqlClient.(sqlclientcontracts.WithDBGetter)
			if !ok {
				return nil, fmt.Errorf("sql client %q does not expose its DB connection", name)
			}

			return CreateMigration(name, dbClient.DB(), cfg, env, logger)
		},
	},
)
