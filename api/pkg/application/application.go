package application

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/application/contracts"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	waiterContracts "frisboo-bank/openapi-generator-service/pkg/waiter/contracts"

	"go.uber.org/dig"
)

var _ contracts.Application = (*application)(nil)

type application struct {
	container *dig.Container
	env       environmentEnum.Environment
	logger    loggerContracts.Logger
	waiter    waiterContracts.Waiter
}

func NewApplication(
	container *dig.Container,
	coreModule containerContracts.Module,
) (contracts.Application, error) {
	validation.AssertNotNil("container", container)
	validation.AssertNotNil("coreModule", coreModule)

	if err := coreModule.Apply(container); err != nil {
		return nil, fmt.Errorf("apply core module: %w", err)
	}

	app := &application{container: container}

	if err := container.Invoke(func(
		env environmentEnum.Environment,
		logger loggerContracts.Logger,
		waiter waiterContracts.Waiter,
	) {
		app.env = env
		app.logger = logger
		app.waiter = waiter
	}); err != nil {
		return nil, fmt.Errorf("invoke dependencies: %w", err)
	}

	if err := coreModule.ApplyHooks(app.container, app.waiter); err != nil {
		return nil, fmt.Errorf("apply hooks: %w", err)
	}

	return app, nil
}

func (a *application) ResolveFunc(function any) {
	if err := a.container.Invoke(function); err != nil {
		a.logger.Errorf("Failed to resolve function: %v", err)
	}
}

func (a *application) Run() error {
	return nil
}

func (a *application) Start() error {
	return a.waiter.Wait()
}

func (a *application) Stop() {
	a.waiter.Cancel()
}

func (a *application) Logger() loggerContracts.Logger {
	return a.logger
}

func (a *application) Environment() environmentEnum.Environment {
	return a.env
}
