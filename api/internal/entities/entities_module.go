package entities

import (
	"fmt"
	"reflect"

	entityv1 "frisboo-bank/openapi-generator-service/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/endpoints"
	listEntitiesHandler "frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/handlers"
	"frisboo-bank/openapi-generator-service/internal/entities/repositories"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	grpcserverContracts "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"

	"go.uber.org/dig"
)

func EntitiesModule() containerContracts.Module {
	return container.NewModule(
		"entities",

		container.Provider(func(params struct {
			dig.In
			SqlClient sqlclientContracts.SQLClientCore `name:"sql-client:main"`
			Logger    loggerContracts.Logger           `name:"logger:main"`
		},
		) (contracts.EntitySQLRepository, error) {
			sqlXClient, ok := params.SqlClient.(sqlclientContracts.SQLXClient)
			if !ok {
				return nil, fmt.Errorf("EntitySQLRepository expected a sqlx db client but %q passed", reflect.TypeOf(params.SqlClient))
			}
			return repositories.NewEntityRepositoryPgx(sqlXClient, params.Logger), nil
		}),

		container.Provider(listEntitiesHandler.NewListEntitiesHandler),

		container.Invoker(func(params struct {
			dig.In
			RpcServer            grpcserverContracts.RPCServer `name:"rpc-server:main"`
			ListEntitiesEndpoint *endpoints.ListEntitiesEndpoint
		},
		) error {
			entityv1.RegisterEntityServiceServer(params.RpcServer, params.ListEntitiesEndpoint)
			return nil
		}),
	)
}
