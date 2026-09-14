package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"frisboo-bank/openapi-generator-service/internal/shared/configurations/generator"
	"frisboo-bank/openapi-generator-service/pkg/application"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type Bootstrap struct {
	configLoader configContracts.ConfigLoader
	env          environmentEnum.Environment
}

func NewBootstrap(
	cfgLoader configContracts.ConfigLoader,
	env environmentEnum.Environment,
) *Bootstrap {
	validation.AssertNotNil("cfgLoader", cfgLoader)
	validation.AssertNotNil("env", env)

	return &Bootstrap{
		configLoader: cfgLoader,
		env:          env,
	}
}

func (b *Bootstrap) Run() error {
	logger := slog.Default()
	logger.Info("bootstrapping application", slog.String("env", string(b.env.String())))

	builder := application.NewApplicationBuilder(b.configLoader, b.env)
	builder.ProvideModule(generator.GeneratorServiceModule(b.configLoader, b.env))

	app, err := builder.Build()
	if err != nil {
		return fmt.Errorf("application stopped with error: %w", err)
	}

	configurator := generator.NewGeneratorServiceConfigurator(app)
	configurator.ConfigureGenerator()

	logger.Info("application starting")

	if err := app.Start(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("application stopped with error: %w", err)
	}

	return nil
}
