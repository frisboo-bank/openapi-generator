package contracts

import (
	"context"

	waiterContracts "frisboo-bank/openapi-generator-service/pkg/waiter/contracts"

	"go.uber.org/dig"
)

type (
	HookResolveResult struct {
		Name    string
		Wait    func(ctx context.Context) error
		Cleanup func(ctx context.Context) error
	}

	Provider struct {
		Fn    any
		Name  string `optional:"true"`
		Group string `optional:"true"`
	}
	Decorator struct{ Fn any }
	Invoker   struct{ Fn any }
	Hook      func(c *dig.Container) (HookResolveResult, error)
)

type Dependency interface {
	IsDependency()
}

func (Provider) IsDependency()  {}
func (Decorator) IsDependency() {}
func (Invoker) IsDependency()   {}
func (Hook) IsDependency()      {}

type Module interface {
	Dependency

	AddProvider(providers ...Provider) Module
	AddDecorator(decorators ...Decorator) Module
	AddInvoker(invokers ...Invoker) Module
	AddHook(hooks ...Hook) Module
	AddModule(modules ...Module) Module
	Apply(container *dig.Container) error
	ApplyHooks(container *dig.Container, waiter waiterContracts.Waiter) error
}
