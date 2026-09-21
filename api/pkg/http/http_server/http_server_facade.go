package httpserver

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

var _ contracts.HTTPServer = (*httpServer)(nil)

type httpServer struct {
	adapter contracts.HTTPServerAdapter
}

func (h *httpServer) AddMiddlewares(middlewares ...any) {
	h.adapter.AddMiddlewares(middlewares...)
}

func (h *httpServer) ListRoutes() []any {
	return h.adapter.ListRoutes()
}

func (h *httpServer) RouteBuilder() contracts.RouteBuilder {
	return h.adapter.RouteBuilder()
}

func (h *httpServer) SetupDefaultMiddlewares() {
	h.adapter.SetupDefaultMiddlewares()
}

func (h *httpServer) Start(ctx context.Context) error {
	return h.adapter.Start(ctx)
}

func (h *httpServer) Stop(ctx context.Context) error {
	return h.adapter.Stop(ctx)
}

func (h *httpServer) Name() string                        { return h.adapter.Name() }
func (h *httpServer) Type() httpservertype.HttpServerType { return h.adapter.Type() }
func (h *httpServer) Logger() loggercontracts.Logger      { return h.adapter.Logger() }
