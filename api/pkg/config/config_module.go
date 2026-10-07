package config

import (
	"frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	config "frisboo-bank/openapi-generator-service/pkg/config/config"
)

func ConfigModule(configLoader contracts.ConfigLoader) containerContracts.Module {
	return container.NewModule(
		"config",

		container.Provider(func() contracts.ConfigLoader {
			return configLoader
		}),

		container.Provider(
			func(loader contracts.ConfigLoader, env environmentEnum.Environment) (*config.AppOptions, error) {
				var cfg config.AppOptions
				if err := loader.LoadKey(env, &cfg, "app"); err != nil {
					return nil, err
				}
				return &cfg, nil
			},
		),
	)
}
