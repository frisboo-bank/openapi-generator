package generator

import (
	"frisboo-bank/openapi-generator-service/internal/entities"
	"frisboo-bank/openapi-generator-service/pkg/cache"
	"frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	sqlclient "frisboo-bank/openapi-generator-service/pkg/database/sql_client"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	httpserver "frisboo-bank/openapi-generator-service/pkg/http/http_server"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	rpcserver "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func GeneratorServiceModule(
	configLoader contracts.ConfigLoader,
	env environmentEnum.Environment,
) containerContracts.Module {
	validation.AssertNotNil("configLoader", configLoader)
	validation.AssertNotNil("env", env)

	return container.NewModule(
		"generator-service",
		httpserver.HTTPServerModule(env, configLoader),
		rpcserver.RPCServerModule(env, configLoader),
		sqlclient.SQLClientModule(env, configLoader),
		cache.CacheModule(env, configLoader),
		mediator.MediatorModule(env, configLoader),

		entities.EntitiesModule(),
	)
}
