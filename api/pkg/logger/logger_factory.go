package logger

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/internal/adapters/noop"
	"frisboo-bank/openapi-generator-service/pkg/logger/internal/adapters/zerolog"
	"frisboo-bank/openapi-generator-service/pkg/logger/models"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func CreateLogger(name string, cfg *models.LoggerOptions, env environmentEnum.Environment) (contracts.Logger, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertValidEnum("env", env)

	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	switch cfg.Type {
	case loggertype.LoggerTypes.ZEROLOG:
		return &logger{zerolog.NewZerologAdapter(name, cfg, env)}, nil
	case loggertype.LoggerTypes.NOOP:
		return &logger{noop.NewNoopAdapter(name)}, nil
	default:
		return nil, syserrors.Newf("no logger of type %q exists", cfg.Type)
	}
}

func CreateNoopLogger(name string, env environmentEnum.Environment) contracts.Logger {
	logger, _ := CreateLogger(name, &models.LoggerOptions{
		Type: loggertype.LoggerTypes.NOOP,
	}, env)
	return logger
}
