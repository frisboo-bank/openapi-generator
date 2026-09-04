package apis

import (
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
)

func ApisModule() containerContracts.Module {
	return container.NewModule(
		"apis",
	)
}
