package generator

import (
	"frisboo-bank/openapi-generator-service/internal/entities"
	entitiesconfigurations "frisboo-bank/openapi-generator-service/internal/entities/configurations"
	"frisboo-bank/openapi-generator-service/pkg/cache"
	"frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	sqlclient "frisboo-bank/openapi-generator-service/pkg/database/sql_client"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	httpserver "frisboo-bank/openapi-generator-service/pkg/http/http_server"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	rpcserver "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func GeneratorServiceModule(
	configLoader contracts.ConfigLoader,
	env environmentenum.Environment,
) containercontracts.Module {
	validation.AssertNotNil("configLoader", configLoader)
	validation.AssertNotNil("env", env)

	return container.NewModule(
		"generator-service",
		httpserver.HTTPServerModule(env, configLoader),
		rpcserver.RPCServerModule(env, configLoader),
		sqlclient.SQLClientModule(env, configLoader),
		cache.CacheModule(env, configLoader),
		mediator.MediatorModule(env, configLoader),
		metrics.MetricsModule(env, configLoader),
		tracer.TracerModule(env, configLoader),

		entities.EntitiesModule(),
		entitiesconfigurations.EntitiesConfigurationsModule(),
	)
}
