package logger

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/adapters/zerolog"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/models"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
)

func NewLogger(name string, cfg *models.LoggerOptions, env environmentEnum.Environment) (contracts.Logger, error) {
	var adapter contracts.Logger

	switch cfg.Type {
	case loggertype.LoggerTypes.ZEROLOG:
		adapter = zerolog.NewZerologAdapter(name, cfg, env)
	default:
		return nil, syserrors.Newf("no logger of type %q exists", cfg.Type)
	}

	return adapter, nil
}
