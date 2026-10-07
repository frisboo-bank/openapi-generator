package rpcserver

import (
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/config"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/internal/adapters/grpc"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/types/rpcservertype"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
)

func CreateRPCServer(
	name string,
	cfg *config.RPCServerOptions,
	logger loggercontracts.Logger,
	env environmentenum.Environment,
) (contracts.RPCServer, error) {
	var adapter contracts.RPCServerAdapter
	var err error

	switch cfg.Type {
	case rpcservertype.RpcServerTypes.GRPC:
		adapter, err = grpc.NewGRPCServer(name, cfg, logger, env), nil
	default:
		err = syserrors.Newf("no rpc-server of type %q exists", cfg.Type)
	}
	if err != nil {
		return nil, err
	}
	return &rpcServer{adapter}, nil
}
