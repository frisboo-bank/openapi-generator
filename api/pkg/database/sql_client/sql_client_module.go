package sqlclient

import (
	"context"
	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"

	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"

	"go.uber.org/dig"
)

type SQLClientDependencies struct {
	dig.In
	Tracer  tracercontracts.Tracer
	Metrics metricscontracts.Metrics
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
					if err := instance.Ping(ctx); err != nil {
						instance.Logger().Fatalf("sql-clients %q failed to access database with error: %v", instance.Name(), err)
					}
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
