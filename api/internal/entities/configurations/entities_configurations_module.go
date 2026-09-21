package configurations

import (
	"context"
	"fmt"
	"net/http"

	"frisboo-bank/openapi-generator-service/internal/entities/configurations/mediator"
	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	httpservercontracts "frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	grpcservercontracts "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/registrar"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func EntitiesConfigurationsModule() containercontracts.Module {
	return container.NewModule(
		"entities_configurations",

		container.Invoker(func(params struct {
			dig.In
			Mediator   mediatorcontracts.Mediator
			Repository contracts.EntityRepository
			Logger     loggercontracts.Logger `name:"logger:main"`
		},
		) error {
			return mediator.ConfigureEntitiesMediator(params.Mediator, params.Repository, params.Logger)
		}),

		container.Invoker(func(params struct {
			dig.In
			RPCServer           grpcservercontracts.RPCServer `name:"rpc-server:main"`
			EntityServiceServer entityv1.EntityServiceServer
		},
		) error {
			svc, ok := params.EntityServiceServer.(registrar.Service)
			if !ok {
				return fmt.Errorf("EntityServiceServer does not implement registrar.Service")
			}
			params.RPCServer.ServiceManager().Register(svc)

			return nil
		}),

		container.Invoker(func(params struct {
			dig.In
			HTTPServer          httpservercontracts.HTTPServer `name:"http-server:main"`
			EntityServiceServer entityv1.EntityServiceServer
		},
		) error {
			conn, err := grpc.NewClient(
				"localhost:9002",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return err
			}

			ctx := context.Background()
			mux := runtime.NewServeMux()

			if err := entityv1.RegisterEntityServiceHandler(ctx, mux, conn); err != nil {
				return err
			}

			rb := params.HTTPServer.RouteBuilder().Root()

			rb.POST("/entities", func(w http.ResponseWriter, r *http.Request) { mux.ServeHTTP(w, r) })
			rb.POST("/entity", func(w http.ResponseWriter, r *http.Request) { mux.ServeHTTP(w, r) })
			rb.GET("/entity/:slug", func(w http.ResponseWriter, r *http.Request) { mux.ServeHTTP(w, r) })
			rb.PATCH("/entity/:slug", func(w http.ResponseWriter, r *http.Request) { mux.ServeHTTP(w, r) })
			rb.DELETE("/entity/:slug", func(w http.ResponseWriter, r *http.Request) { mux.ServeHTTP(w, r) })

			return nil
		}),
	)
}
