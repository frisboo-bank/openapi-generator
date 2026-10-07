package logger

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/internal/adapters/noop"
)

func CreateNoopLogger(name string, env environmentEnum.Environment) contracts.Logger {
	return &logger{adapter: noop.NewNoopAdapter(name)}
}
