package mediator

import (
	"fmt"
	"log"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator/models"
)

func MediatorModule(
	env environmentEnum.Environment,
	configLoader configContracts.ConfigLoader,
) containerContracts.Module {
	var cfg *models.MediatorOptions
	if err := configLoader.LoadKey(env, &cfg, "mediator"); err != nil {
		log.Fatalf("Failed to build mediator module with error: %v", err)
	}

	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Failed to build mediator module with error: %v", err)
	}

	mod := container.NewModule(
		"mediator",
		container.Provider(func() *models.MediatorOptions {
			return cfg
		}),
	)

	mod.AddProvider(container.Provider(func(
		loggers map[string]loggerContracts.Logger,
	) (contracts.Mediator, error) {
		loggerName := cfg.Logger
		logger, ok := loggers[loggerName]
		if !ok {
			return nil, fmt.Errorf("logger %q not found", loggerName)
		}
		return NewMediator(cfg, logger)
	}))

	return mod
}
