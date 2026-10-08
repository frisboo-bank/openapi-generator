package migration

import (
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	migrationinternal "frisboo-bank/openapi-generator-service/pkg/database/migration/internal"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	SQLClients map[string]sqlclientcontracts.SQLClientCore
	Tracers    map[string]tracercontracts.Tracer
	Metrics    map[string]metricscontracts.Metrics
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
			sqlClient, ok := extra.SQLClients[cfg.DBClient]
			if !ok {
				return nil, fmt.Errorf("sql client %q not found for migration %q", cfg.DBClient, name)
			}
			dbClient, ok := sqlClient.(sqlclientcontracts.WithDBGetter)
			if !ok {
				return nil, fmt.Errorf("sql client %q does not expose its DB connection", name)
			}

			var tracerInstance tracercontracts.Tracer
			if cfg.EnableTracing {
				tracerInstance, ok = extra.Tracers[cfg.Tracer]
				if !ok {
					return nil, fmt.Errorf("tracer %q not found for migration %q", cfg.Tracer, name)
				}
			}

			var metricsInstance metricscontracts.Metrics
			if cfg.EnableMetrics {
				metricsInstance, ok = extra.Metrics[cfg.Metrics]
				if !ok {
					return nil, fmt.Errorf("metrics %q not found for migration %q", cfg.Metrics, name)
				}
			}

			return migrationinternal.CreateMigration(name, dbClient.DB(), cfg, env, logger, tracerInstance, metricsInstance)
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

	return migrationinternal.CreateMigrationForTesting(
		name,
		db,
		&config.MigrationOptions{
			MigrationsDir: migrationDir,
			Type:          migrationtype.MigrationTypes.GOOSE,
		},
		nil,
	)
}
