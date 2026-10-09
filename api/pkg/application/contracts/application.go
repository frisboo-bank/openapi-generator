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
	Stop()
}
