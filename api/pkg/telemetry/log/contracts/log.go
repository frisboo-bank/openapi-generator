package contracts

import (
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	logtype "frisboo-bank/openapi-generator-service/pkg/telemetry/log/models/enums/log_type"
)

type (
	Log interface {
		LogAdapter
	}

	LogAdapter interface {
		Name() string
		Type() logtype.LogType
		Logger() loggercontracts.Logger
	}
)
