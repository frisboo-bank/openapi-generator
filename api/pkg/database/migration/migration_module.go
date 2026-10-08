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
	"frisboo-bank/openapi-generator-service/pkg/logger"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"
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
		Name:      "database.migration",
		ConfigKey: "database.migration",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			types.MigrationEnumsDecodeHook(),
		},
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

			// Telemetry is opt-in via enableTelemetry. When off, resolve
			// noop instances so the decorator is a no-op and the adapter runs
			// uninstrumented — matching the sql_client pattern.
			var tracerInstance tracercontracts.Tracer
			if !cfg.EnableTelemetry {
				var err error
				tracerInstance, err = tracer.CreateNoopTracer(name, logger)
				if err != nil {
					return nil, fmt.Errorf("failed to create tracer for migration %q: %w", name, err)
				}
			} else {
				var ok bool
				tracerInstance, ok = extra.Tracers[cfg.Tracer]
				if !ok {
					return nil, fmt.Errorf("tracer %q not found for migration %q", cfg.Tracer, name)
				}
			}

			var metricsInstance metricscontracts.Metrics
			if !cfg.EnableTelemetry {
				var err error
				metricsInstance, err = metrics.CreateNoopMetrics(name, logger)
				if err != nil {
					return nil, fmt.Errorf("failed to create metrics for migration %q: %w", name, err)
				}
			} else {
				var ok bool
				metricsInstance, ok = extra.Metrics[cfg.Metrics]
				if !ok {
					return nil, fmt.Errorf("metrics %q not found for migration %q", cfg.Metrics, name)
				}
			}

			return migrationinternal.CreateMigration(name, dbClient.DB(), cfg, env, logger, tracerInstance, metricsInstance)
		},
	},
)

// CreateMigrationForTests is the test-only entry point: it skips the
// development environment gate so migrations can run under the TESTING
// environment.
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
		logger.CreateNoopLogger("test", environmentenum.Environments.TESTING),
	)
}
