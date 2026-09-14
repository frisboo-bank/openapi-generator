package generator

import (
	"fmt"
	"net/http"

	"frisboo-bank/openapi-generator-service/pkg/application/contracts"
	configModels "frisboo-bank/openapi-generator-service/pkg/config/models"
	"frisboo-bank/openapi-generator-service/pkg/container"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	httpServerContracts "frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/endpoints"
	rpcServerContracts "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	grpcservices "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/services/grpc"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"go.uber.org/dig"
)

type GeneratorServiceConfigurator struct {
	app                        contracts.Application
	infrastructureConfigurator *InfrastructureConfigurator
}

func NewGeneratorServiceConfigurator(app contracts.Application) *GeneratorServiceConfigurator {
	validation.AssertNotNil("app", app)

	infraConfigurator := NewInfrastructureConfigurator(app)

	return &GeneratorServiceConfigurator{
		app:                        app,
		infrastructureConfigurator: infraConfigurator,
	}
}

func (c *GeneratorServiceConfigurator) ConfigureGenerator() {
	c.infrastructureConfigurator.ConfigureInfrastructures()
	c.mapGeneratorEndpoints()
	c.mapGeneratorRPCServices()
}

func (c *GeneratorServiceConfigurator) mapGeneratorEndpoints() {
	c.app.ResolveFunc(container.Invoker(
		func(params struct {
			dig.In
			Server httpServerContracts.HTTPServer `name:"http-server:main"`
			Cfg    *configModels.AppOptions
			Env    environmentEnum.Environment
		},
		) {
			params.Server.SetupDefaultMiddlewares()

			rb := params.Server.RouteBuilder().Root()

			if params.Env.IsDevelopment() {
				endpoints.NewChromeDevtoolEndpoint(params.Env, params.Server.Logger(), rb).MapEndpoint()
			}

			rb.GET("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, "%s is running", params.Cfg.Name)
			})
		},
	).Fn)
}

func (c *GeneratorServiceConfigurator) mapGeneratorRPCServices() {
	c.app.ResolveFunc(container.Invoker(
		func(params struct {
			dig.In
			Server rpcServerContracts.RPCServer `name:"rpc-server:main"`
		},
		) {
			sm := params.Server.ServiceManager()

			params.Server.SetupDefaultMiddlewares()
			sm.Register(grpcservices.NewGRPCHealthService("main", params.Server.Logger()))
		},
	).Fn)
}
