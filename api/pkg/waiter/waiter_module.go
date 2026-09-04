package waiter

import (
	"fmt"
	"log"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/models"
)

func WaiterModule(
	env environmentEnum.Environment,
	configLoader configContracts.ConfigLoader,
) containerContracts.Module {
	var cfg *models.WaiterOptions
	if err := configLoader.LoadKey(env, &cfg, "waiter"); err != nil {
		log.Fatalf("Failed to build waiter module with error: %v", err)
	}

	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Failed to build waiter module with error: %v", err)
	}

	mod := container.NewModule(
		"waiter",
		container.Provider(func() *models.WaiterOptions {
			return cfg
		}),
	)

	mod.AddProvider(container.Provider(func(
		loggers map[string]loggerContracts.Logger,
	) (contracts.Waiter, error) {
		loggerName := cfg.Logger
		logger, ok := loggers[loggerName]
		if !ok {
			return nil, fmt.Errorf("logger %q not found", loggerName)
		}
		return NewWaiter(cfg, logger)
	}))

	return mod
}
