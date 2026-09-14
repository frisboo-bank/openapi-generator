package contracts

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

type ApplicationBuilder interface {
	ProvideModule(module ...containerContracts.Module)
	ProvideProvider(providers ...any)
	ProvideDecorator(decorators ...any)
	Build() (Application, error)
	ConfigLoader() configContracts.ConfigLoader
	Environment() environmentEnum.Environment
}
