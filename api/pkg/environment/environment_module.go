package environment

import (
	configloadercontracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

func EnvironmentModule(env environmentenum.Environment, configLoader configloadercontracts.ConfigLoader) containercontracts.Module {
	configLoader.RegisterDecodeHookFunc(environmentenum.EnvironmentEnumsDecodeHook())
	return container.NewModule(
		"environment",
		container.Provider(func() environmentenum.Environment { return env }),
	)
}
