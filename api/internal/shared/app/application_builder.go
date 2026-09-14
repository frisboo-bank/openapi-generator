package app

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/shared/configurations/generator"
	"frisboo-bank/openapi-generator-service/pkg/application"
	"frisboo-bank/openapi-generator-service/pkg/application/contracts"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type GeneratorApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewGeneratorApplicationBuilder(
	env environmentEnum.Environment,
	configLoader configContracts.ConfigLoader,
) *GeneratorApplicationBuilder {
	validation.AssertNotNil("env", env)
	validation.AssertNotNil("configLoader", configLoader)

	return &GeneratorApplicationBuilder{
		ApplicationBuilder: application.NewApplicationBuilder(configLoader, env),
	}
}

func (b *GeneratorApplicationBuilder) Build(ctx context.Context) (contracts.Application, error) {
	b.ProvideModule(generator.GeneratorServiceModule(
		b.ConfigLoader(),
		b.Environment(),
	))
	return b.ApplicationBuilder.Build()
}
