package migration

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/models"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"

	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	SQLClients map[string]sqlclientcontracts.SQLClientCore
}

var MigrationModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.MigrationOptions, contracts.Migration, MigrationDependencies]{
		Name:       "migration",
		ConfigKey:  "migration",
		ProviderFn: CreateMigration,
	},
)
