package environment

import (
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

func EnvironmentModule(env environmentEnum.Environment) containerContracts.Module {
	return container.NewModule(
		"environment",
		container.Provider(func() environmentEnum.Environment { return env }),
	)
}
