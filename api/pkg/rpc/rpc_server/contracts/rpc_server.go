package contracts

import (
	"context"

	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/registrar"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/types/rpcservertype"
)

type (
	RPCServer interface {
		RPCServerAdapter
	}

	RPCServerAdapter interface {
		AddMiddlewares(middlewares ...any)
		ListServices() []any
		Logger() loggerContracts.Logger
		Name() string
		ServiceManager() *registrar.ServiceManager
		SetupDefaultMiddlewares()
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		Type() rpcservertype.RpcServerType
	}
)
