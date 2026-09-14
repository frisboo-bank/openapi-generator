package rpcserver

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models"

	"go.uber.org/dig"
)

type RPCServerDependencies struct {
	dig.In
}

var RPCServerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*models.RPCServerOptions, contracts.RPCServer, RPCServerDependencies]{
		Name:      "rpc-server",
		ConfigKey: "rpc-servers",
		ProviderFn: func(
			name string,
			cfg *models.RPCServerOptions,
			env environmentEnum.Environment,
			logger loggerContracts.Logger,
			extra RPCServerDependencies,
		) (contracts.RPCServer, error) {
			return CreateRPCServer(name, cfg, logger, env)
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
