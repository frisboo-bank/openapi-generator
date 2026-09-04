package generator

import (
	"frisboo-bank/openapi-generator-service/pkg/application/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(app contracts.Application) *InfrastructureConfigurator {
	validation.AssertNotNil("app", app)

	return &InfrastructureConfigurator{
		Application: app,
	}
}

func (ic *InfrastructureConfigurator) ConfigureInfrastructures() {
	ic.ResolveFunc(
		func() {},
	)
}
