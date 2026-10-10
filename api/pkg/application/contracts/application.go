package contracts

import (
	"context"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type Application interface {
	Environment() environmentEnum.Environment
	Logger() loggerContracts.Logger
	ResolveFunc(function any)
	Start(ctx context.Context) error

	// Stop requests graceful shutdown and blocks until it completes.
	// Pass a context with a deadline to force shutdown if cleanup does not
	// finish in time.
	Stop(ctx context.Context) error
}
