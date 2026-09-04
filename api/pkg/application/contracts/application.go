package contracts

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type Application interface {
	Environment() environmentEnum.Environment
	Logger() loggerContracts.Logger
	ResolveFunc(function any)
	Run() error
	Start() error
	Stop()
}
