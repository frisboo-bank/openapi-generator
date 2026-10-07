package logger

import (
	"fmt"

	configloadercontracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/config"
	loggerinternal "frisboo-bank/openapi-generator-service/pkg/logger/internal"
	loggerenums "frisboo-bank/openapi-generator-service/pkg/logger/types"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func LoggerModule(env environmentenum.Environment, configLoader configloadercontracts.ConfigLoader) containercontracts.Module {
	configLoader.RegisterDecodeHookFunc(loggerenums.LoggerEnumsDecodeHook())

	var cfgMap map[string]*config.LoggerOptions
	if err := configLoader.LoadKey(env, &cfgMap, "loggers"); err != nil {
		panic(fmt.Sprintf("failed to load logger config: %v", err))
	}

	mod := container.NewModule(
		"logger",
		container.Provider(func() map[string]*config.LoggerOptions {
			return cfgMap
		}),
	)

	loggersMap := make(map[string]contracts.Logger, len(cfgMap))

	for name, cfg := range cfgMap {
		validation.AssertNotNil("cfg", cfg)

		logger, err := loggerinternal.CreateLogger(name, cfg, env)
		if err != nil {
			panic(fmt.Sprintf("failed to create logger %q: %v", name, err))
		}

		mod.AddProvider(containercontracts.Provider{
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
