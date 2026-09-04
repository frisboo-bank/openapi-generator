package application

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/application/contracts"
	"frisboo-bank/openapi-generator-service/pkg/config"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	configModels "frisboo-bank/openapi-generator-service/pkg/config/models"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	"frisboo-bank/openapi-generator-service/pkg/environment"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"frisboo-bank/openapi-generator-service/pkg/waiter"

	"go.uber.org/dig"
)

var _ contracts.ApplicationBuilder = (*applicationBuilder)(nil)

type applicationBuilder struct {
	container    *dig.Container
	configLoader configContracts.ConfigLoader
	modules      []containerContracts.Module
	providers    []containerContracts.Provider
	decorators   []containerContracts.Decorator
	environment  environmentEnum.Environment
}

func NewApplicationBuilder(
	configLoader configContracts.ConfigLoader,
	env environmentEnum.Environment,
) contracts.ApplicationBuilder {
	validation.AssertNotNil("configLoader", configLoader)
	validation.AssertNotNil("env", env)

	return &applicationBuilder{
		container:    dig.New(),
		configLoader: configLoader,
		environment:  env,
	}
}

func (a *applicationBuilder) Build() (contracts.Application, error) {
	coreModule := container.NewModule(
		"core",
		environment.EnvironmentModule(a.environment),
		config.ConfigModule(a.configLoader),
		logger.LoggerModule(a.environment, a.configLoader),
		waiter.WaiterModule(a.environment, a.configLoader),

		container.Provider(func(
			loggers map[string]loggerContracts.Logger,
			cfg *configModels.AppOptions,
		) (loggerContracts.Logger, error) {
			loggerName := cfg.Logger
			for _, l := range loggers {
				if l.Name() == loggerName {
					return l, nil
				}
			}

			return nil, fmt.Errorf("logger %q not found for app", loggerName)
		}),
	)

	coreModule.
		AddProvider(a.providers...).
		AddDecorator(a.decorators...).
		AddModule(a.modules...)

	return NewApplication(
		a.container,
		coreModule,
	)
}

func (a *applicationBuilder) ConfigLoader() configContracts.ConfigLoader {
	return a.configLoader
}

func (a *applicationBuilder) Environment() environmentEnum.Environment {
	return a.environment
}

func (a *applicationBuilder) ProvideDecorator(decorators ...any) {
	for _, d := range decorators {
		a.decorators = append(a.decorators, container.Decorator(d))
	}
}

func (a *applicationBuilder) ProvideModule(modules ...containerContracts.Module) {
	a.modules = append(a.modules, modules...)
}

func (a *applicationBuilder) ProvideProvider(providers ...any) {
	for _, p := range providers {
		a.providers = append(a.providers, container.Provider(p))
	}
}
