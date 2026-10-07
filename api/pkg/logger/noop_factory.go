package logger

import (
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/config"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	loggerinternal "frisboo-bank/openapi-generator-service/pkg/logger/internal"
	"frisboo-bank/openapi-generator-service/pkg/logger/types/loggertype"
)

func CreateNoopLogger(name string, env environmentenum.Environment) (contracts.Logger, error) {
	return loggerinternal.CreateLogger(name, &config.LoggerOptions{Type: loggertype.LoggerTypes.NOOP}, env)
}
