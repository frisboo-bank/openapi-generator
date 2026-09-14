package outbox

import (
	configLoaderContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

func NewOutBoxModule(
	env environmentEnum.Environment,
	configLoader configLoaderContracts.ConfigLoader,
) containerContracts.Module {
	return container.NewModule("outbox")
}
