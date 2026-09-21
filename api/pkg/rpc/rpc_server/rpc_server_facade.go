package rpcserver

import (
	"context"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/registrar"
)

var _ contracts.RPCServer = (*rpcServer)(nil)

type rpcServer struct {
	adapter contracts.RPCServerAdapter
}

func (r *rpcServer) AddMiddlewares(middlewares ...any) {
	r.adapter.AddMiddlewares(middlewares...)
}

func (r *rpcServer) ListServices() []any {
	return r.adapter.ListServices()
}

func (r *rpcServer) SetupDefaultMiddlewares() {
	r.adapter.SetupDefaultMiddlewares()
}

func (r *rpcServer) Start(ctx context.Context) error {
	return r.adapter.Start(ctx)
}

func (r *rpcServer) Stop(ctx context.Context) error {
	return r.adapter.Stop(ctx)
}

func (r *rpcServer) Name() string                              { return r.adapter.Name() }
func (r *rpcServer) Type() rpcservertype.RpcServerType         { return r.adapter.Type() }
func (r *rpcServer) Logger() loggercontracts.Logger            { return r.adapter.Logger() }
func (r *rpcServer) ServiceManager() *registrar.ServiceManager { return r.adapter.ServiceManager() }
