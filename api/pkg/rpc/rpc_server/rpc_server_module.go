package rpcserver

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/config"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	rpcserverinternal "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/internal"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/types"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type RPCServerDependencies struct {
	dig.In
	Loggers     module.DependenciesMap[loggercontracts.Logger]
	Environment environmentenum.Environment
}

var RPCServerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.RPCServerOptions, contracts.RPCServer, RPCServerDependencies]{
		Name:      "rpc-server",
		ConfigKey: "rpc-servers",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			types.RPCServerEnumsDecodeHook(),
		},
		ProviderFn: func(
			ctx context.Context,
			name string,
			cfg *config.RPCServerOptions,
			dependencies RPCServerDependencies,
		) (contracts.RPCServer, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return rpcserverinternal.CreateRPCServer(name, cfg, loggerInstance, dependencies.Environment)
		},
		HookFn: func(ctx context.Context, name string, instance contracts.RPCServer) containercontracts.HookResolveResult {
			return containercontracts.HookResolveResult{
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
