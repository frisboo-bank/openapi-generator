package contracts

import (
	"context"

	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"
)

type (
	RPCServer interface {
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		Name() string
		Type() rpcservertype.RpcServerType
		Logger() loggerContracts.Logger
	}
)
