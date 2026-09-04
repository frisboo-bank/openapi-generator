package rpcserver

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/adapters/grpc"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models"
	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"

	"go.uber.org/dig"
)

type RPCServerDependencies struct {
	dig.In
}

var RPCServerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.RPCServerOptions, contracts.RPCServer, RPCServerDependencies]{
		Name:      "rpc-server",
		ConfigKey: "rpc-servers",
		ProviderFn: func(name string,
			cfg *models.RPCServerOptions,
			env environmentEnum.Environment,
			logger loggerContracts.Logger,
			extra RPCServerDependencies,
		) (contracts.RPCServer, error) {
			switch cfg.Type {
			case rpcservertype.RpcServerTypes.GRPC:
				return grpc.NewGRPCServer(name, cfg, logger, env), nil
			default:
				return nil, syserrors.Newf("no rpc-server of type %q exists", cfg.Type)
			}
		},
		HookFn: func(name string, instance contracts.RPCServer) containerContracts.HookResolveResult {
			return containerContracts.HookResolveResult{
				Name: "rpc-server:" + name,
				Wait: func(ctx context.Context) error {
					return instance.Start(ctx)
				},
				Cleanup: func(ctx context.Context) error {
					if err := instance.Stop(ctx); err != nil {
						instance.Logger().Errorf("rpc-server %q shutdown failed: %v", name, err)
						return err
					}

					instance.Logger().Infof("rpc-server %q shut down gracefully", name)

					return nil
				},
			}
		},
	},
)
