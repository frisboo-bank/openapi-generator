package contracts

import (
	"context"

	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type (
	HTTPServer interface {
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		ListRoutes() []any
		Name() string
		Type() httpservertype.HttpServerType
		RouteBuilder() RouteBuilder
		Logger() loggerContracts.Logger
	}
)
