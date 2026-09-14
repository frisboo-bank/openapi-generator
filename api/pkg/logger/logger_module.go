package logger

import (
	"fmt"

	configLoaderContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/models"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func LoggerModule(env environmentEnum.Environment, configLoader configLoaderContracts.ConfigLoader) containerContracts.Module {
	var cfgMap map[string]*models.LoggerOptions
	if err := configLoader.LoadKey(env, &cfgMap, "loggers"); err != nil {
		panic(fmt.Sprintf("failed to load logger config: %v", err))
	}

	mod := container.NewModule(
		"logger",
		container.Provider(func() map[string]*models.LoggerOptions {
			return cfgMap
		}),
	)

	loggersMap := make(map[string]contracts.Logger, len(cfgMap))

	for name, cfg := range cfgMap {
		validation.AssertNotNil("cfg", cfg)

		logger, err := CreateLogger(name, cfg, env)
		if err != nil {
			panic(fmt.Sprintf("failed to create logger %q: %v", name, err))
		}

		mod.AddProvider(containerContracts.Provider{
			Fn:   func() contracts.Logger { return logger },
			Name: "logger:" + name,
		})

		loggersMap[name] = logger
	}

	mod.AddProvider(container.Provider(
		func() map[string]contracts.Logger { return loggersMap },
	))

	return mod
}
