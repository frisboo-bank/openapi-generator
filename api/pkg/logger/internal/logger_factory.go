package logger

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/config"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/internal/adapters/zerolog"
	"frisboo-bank/openapi-generator-service/pkg/logger/types/loggertype"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func CreateLogger(
	name string,
	cfg *config.LoggerOptions,
	env environmentEnum.Environment,
) (contracts.Logger, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertValidEnum("env", env)

	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	var logger contracts.Logger
	var err error

	switch cfg.Type {
	case loggertype.LoggerTypes.ZEROLOG:
		logger = zerolog.NewZerologAdapter(name, cfg, env)
	default:
		return nil, syserrors.Newf("no logger of type %q exists", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	return logger, nil
}
