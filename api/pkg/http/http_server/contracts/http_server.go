package contracts

import (
	"context"

	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type (
	HTTPServer interface {
		HTTPServerAdapter
	}

	HTTPServerAdapter interface {
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		ListRoutes() []any
		RouteBuilder() RouteBuilder
		Name() string
		Type() httpservertype.HttpServerType
		Logger() loggerContracts.Logger
	}
)
