package migration

import (
	"context"
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
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MigrationDependencies struct {
	dig.In
	Loggers     module.DependenciesMap[loggercontracts.Logger]
	SQLClients  module.DependenciesMap[sqlclientcontracts.SQLClientCore]
	Tracers     module.DependenciesMap[tracercontracts.Tracer]
	Metrics     module.DependenciesMap[metricscontracts.Metrics]
	Environment environmentenum.Environment
}

var MigrationModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.MigrationOptions, contracts.Migration, MigrationDependencies]{
		Name:             "database.migration",
		ConfigKey:        "database.migration",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.MigrationEnumsDecodeHook()},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.MigrationOptions,
			dependencies MigrationDependencies,
		) (contracts.Migration, error) {
			validation.AssertNotNil("ctx", ctx)
			validation.AssertNotEmpty("name", name)
			validation.AssertNotNil("cfg", cfg)
			validation.AssertNotNil("dependencies", dependencies)

			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			sqlClient, err := dependencies.SQLClients.Get(cfg.DBClient)
			if err != nil {
				return nil, err
			}
			dbClient, ok := sqlClient.(sqlclientcontracts.WithDBGetter)
			if !ok {
				return nil, fmt.Errorf("sql client %q does not expose its DB connection", cfg.DBClient)
			}

			var tracerInstance tracercontracts.Tracer
			if cfg.EnableTracing {
				tracerInstance, err = dependencies.Tracers.Get(cfg.Tracer)
				if err != nil {
					return nil, err
				}
			}

			var metricsInstance metricscontracts.Metrics
			if cfg.EnableMetrics {
				metricsInstance, err = dependencies.Metrics.Get(cfg.Metrics)
				if err != nil {
					return nil, err
				}
			}

			return migrationinternal.CreateMigration(
				name,
				dbClient.DB(),
				cfg,
				dependencies.Environment,
				loggerInstance,
				tracerInstance,
				metricsInstance,
			)
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
