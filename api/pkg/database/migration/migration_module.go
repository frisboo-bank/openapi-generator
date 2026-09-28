package migration

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/models"
	migrationenums "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	SQLClients map[string]sqlclientcontracts.SQLClientCore
}

var MigrationModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.MigrationOptions, contracts.Migration, MigrationDependencies]{
		Name:      "migration",
		ConfigKey: "migration",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			migrationenums.MigrationEnumsDecodeHook(),
		},
		ProviderFn: CreateMigration,
	},
)
