package sqlclient

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"

	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	sqlclientenums "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type SQLClientDependencies struct {
	dig.In
	Tracer  tracercontracts.Tracer   `name:"telemetry.tracer:main"`
	Metrics metricscontracts.Metrics `name:"telemetry.metrics:main"`
}

var SQLClientModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.SQLClientOptions, contracts.SQLClientCore, SQLClientDependencies]{
		Name:      "database.sql-client",
		ConfigKey: "database.sql-clients",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			sqlclientenums.SQLClientEnumsDecodeHook(),
		},
		ProviderFn: func(name string, cfg *models.SQLClientOptions, env environmentenum.Environment, logger loggercontracts.Logger, extra SQLClientDependencies) (contracts.SQLClientCore, error) {
			return CreateSQLClient(name, cfg, logger, extra.Tracer, extra.Metrics)
		},
		HookFn: func(name string, instance contracts.SQLClientCore) containercontracts.HookResolveResult {
			return containercontracts.HookResolveResult{
				Name: "database.sql-client:" + name,
				Wait: func(ctx context.Context) error {
					if err := instance.Ping(ctx); err != nil {
						instance.Logger().Errorf("sql-clients %q failed to access database: %v", name, err)
						return fmt.Errorf("sql-client %q ping: %w", name, err)
					}

					<-ctx.Done()
					return nil
				},
				Cleanup: func(ctx context.Context) error {
					if err := instance.Close(ctx); err != nil {
						instance.Logger().Errorf("sql-clients %q close failed with error: %v", name, err)
						return fmt.Errorf("sql-client %q close: %w", name, err)
					}

					instance.Logger().Infof("sql-clients: %q shutdown successfully", name)
					return nil
				},
			}
		},
	},
)
