package outbox

import (
	configLoaderContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

func NewOutBoxModule(
	env environmentEnum.Environment,
	configLoader configLoaderContracts.ConfigLoader,
) containerContracts.Module {
}
