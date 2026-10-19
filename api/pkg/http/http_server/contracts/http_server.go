package contracts

import (
	"context"

	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/types/httpservertype"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/types/hookstatus"
)

type HTTPServer interface {
	SetupDefaultMiddlewares()
	AddMiddlewares(middlewares ...any)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status() hookstatus.HookStatus
	ListRoutes() []any
	RouteBuilder() RouteBuilder
	Name() string
	Type() httpservertype.HttpServerType
	Logger() loggerContracts.Logger
}
