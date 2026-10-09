package contracts

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/types/httpservertype"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type HTTPServer interface {
	SetupDefaultMiddlewares()
	AddMiddlewares(middlewares ...any)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	ListRoutes() []any
	RouteBuilder() RouteBuilder
	Name() string
	Type() httpservertype.HttpServerType
	Logger() loggercontracts.Logger
}
