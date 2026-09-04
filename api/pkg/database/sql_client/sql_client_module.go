package sqlclient

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/adapters/postgres/pg"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/adapters/postgres/pgx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"go.uber.org/dig"
)

type SQLClientDependencies struct {
	dig.In
}

var SQLClientModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.SQLClientOptions, contracts.SQLClientCore, SQLClientDependencies]{
		Name:      "sql-client",
		ConfigKey: "sql-clients",
		ProviderFn: func(
			name string,
			cfg *models.SQLClientOptions,
			env environmentEnum.Environment,
			logger loggerContracts.Logger,
			extra SQLClientDependencies,
		) (contracts.SQLClientCore, error) {
			switch cfg.Type {
			case sqlclienttype.SqlClientTypes.POSTGRES:
				return pg.NewPostgresSQLClientAdapter(name, cfg, logger)
			case sqlclienttype.SqlClientTypes.POSTGRESX:
				return pgx.NewPostgresSQLXClientAdapter(name, cfg, logger)
			default:
				return nil, fmt.Errorf("unsupported SQLClient type: %v", cfg.Type)
			}
		},
		HookFn: func(name string, instance contracts.SQLClientCore) containerContracts.HookResolveResult {
			return containerContracts.HookResolveResult{
				Name: "sql-client:" + name,
				Wait: func(ctx context.Context) error {
					go func() {
						if err := instance.Ping(ctx); err != nil {
							instance.Logger().
								Fatalf("sql-clients %q failed to access database with error: %v", instance.Name(), err)
						}
					}()

					// Wait until the context is canceled (signal/timeout).
					<-ctx.Done()

					return nil
				},
				Cleanup: func(ctx context.Context) error {
					if err := instance.Close(); err != nil {
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
