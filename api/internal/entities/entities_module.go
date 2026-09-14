package entities

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	createentityendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/endpoints"
	createentityqueries "frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/queries"
	deleteentityendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/endpoints"
	deleteentityqueries "frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/queries"
	getentitybyslugendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/endpoints"
	getentitybyslugqueries "frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/queries"
	listentitiesendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/endpoints"
	listentitiesqueries "frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/queries"
	updateentityendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/endpoints"
	updateentityqueries "frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/queries"
	"frisboo-bank/openapi-generator-service/internal/entities/repositories"
	"frisboo-bank/openapi-generator-service/internal/entities/services"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"go.uber.org/dig"
)

func EntitiesModule() containercontracts.Module {
	return container.NewModule(
		"entities",

		container.Provider(func(params struct {
			dig.In
			SQLClient sqlclientcontracts.SQLClientCore `name:"sql-client:main"`
			Logger    loggercontracts.Logger           `name:"logger:main"`
		},
		) (contracts.EntityRepository, error) {
			sqlXClient, ok := params.SQLClient.(sqlclientcontracts.SQLXClientAdapter)
			if !ok {
				return nil, fmt.Errorf("expected SQLXClient, got %T", params.SQLClient)
			}
			return repositories.NewEntityRepositoryPgx(sqlXClient, params.Logger), nil
		}),

		container.Provider(createentityqueries.NewCreateEntityHandler),
		container.Provider(createentityendpoints.NewCreateEntityEndpoint),
		container.Provider(deleteentityqueries.NewDeleteEntityHandler),
		container.Provider(deleteentityendpoints.NewDeleteEntityEndpoint),
		container.Provider(getentitybyslugqueries.NewGetEntityBySlugHandler),
		container.Provider(getentitybyslugendpoints.NewGetEntityBySlugEndpoint),
		container.Provider(listentitiesqueries.NewListEntitiesHandler),
		container.Provider(listentitiesendpoints.NewListEntitiesEndpoint),
		container.Provider(updateentityqueries.NewUpdateEntityHandler),
		container.Provider(updateentityendpoints.NewUpdateEntityEndpoint),

		container.Provider(services.NewEntityGRPCServerService),
	)
}
