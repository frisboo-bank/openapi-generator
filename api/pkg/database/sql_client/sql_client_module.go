package sqlclient

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	"go.uber.org/dig"
)

type SQLClientDependencies struct {
	dig.In
	Logger  loggercontracts.Logger   `optional:"true"`
	Tracer  tracercontracts.Tracer   `optional:"true"`
	Metrics metricscontracts.Metrics `optional:"true"`
}

var SQLClientModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.SQLClientOptions, contracts.SQLClientCore, SQLClientDependencies]{
		Name:      "sql-client",
		ConfigKey: "sql-clients",
		ProviderFn: func(name string, cfg *models.SQLClientOptions, env environmentEnum.Environment, logger loggercontracts.Logger, extra SQLClientDependencies) (contracts.SQLClientCore, error) {
			return CreateSQLClient(name, cfg, logger, extra.Tracer, extra.Metrics)
		},
		HookFn: func(name string, instance contracts.SQLClientCore) containercontracts.HookResolveResult {
			return containercontracts.HookResolveResult{
				Name: "sql-client:" + name,
				Wait: func(ctx context.Context) error {
					// go func() {
					if err := instance.Ping(ctx); err != nil {
						instance.Logger().Fatalf("sql-clients %q failed to access database with error: %v", instance.Name(), err)
					}
					// }()

					// Wait until the context is canceled (signal/timeout).
					<-ctx.Done()

					return nil
				},
				Cleanup: func(ctx context.Context) error {
					if err := instance.Close(ctx); err != nil {
						instance.Logger().Errorf("sql-clients %q close failed with error: %v", name, err)
						return err
					}

					instance.Logger().Infof("sql-clients: %q shutdown successfully", name)
					return nil
				},
			}
		},
	},
)
