package migration

import (
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	migrationinternal "frisboo-bank/openapi-generator-service/pkg/database/migration/internal"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	SQLClients map[string]sqlclientcontracts.SQLClientCore
}

var MigrationModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.MigrationOptions, contracts.Migration, MigrationDependencies]{
		Name:             "database.migration",
		ConfigKey:        "database.migration",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.MigrationEnumsDecodeHook()},
		ProviderFn: func(
			name string,
			cfg *config.MigrationOptions,
			env environmentenum.Environment,
			logger loggercontracts.Logger,
			extra MigrationDependencies,
		) (contracts.Migration, error) {
			return migrationinternal.CreateMigration(name, cfg, env, logger, extra)
		},
	},
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

	return migrationinternal.CreateMigration(
		name,
		&config.MigrationOptions{
			MigrationsDir: migrationDir,
		},
		env,
		logger.CreateNoopLogger("test", environmentenum.Environments.TESTING),
		MigrationDependencies{},
	)
}
