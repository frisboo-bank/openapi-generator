package container

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/container/contracts"
	waiterContracts "frisboo-bank/openapi-generator-service/pkg/waiter/contracts"

	"go.uber.org/dig"
)

var _ contracts.Module = (*module)(nil)

type module struct {
	name       string
	providers  []contracts.Provider
	decorators []contracts.Decorator
	hooks      []contracts.Hook
	invokers   []contracts.Invoker
	modules    []contracts.Module
}

func Provider(fn any) contracts.Provider {
	return contracts.Provider{Fn: fn}
}

func Decorator(fn any) contracts.Decorator {
	return contracts.Decorator{Fn: fn}
}

func Invoker(fn any) contracts.Invoker {
	return contracts.Invoker{Fn: fn}
}

func NewModule(name string, deps ...contracts.Dependency) contracts.Module {
	m := &module{name: name}

	for _, d := range deps {
		switch v := d.(type) {
		case contracts.Provider:
			m.providers = append(m.providers, v)
		case contracts.Decorator:
			m.decorators = append(m.decorators, v)
		case contracts.Invoker:
			m.invokers = append(m.invokers, v)
		case contracts.Hook:
			m.hooks = append(m.hooks, v)
		case contracts.Module:
			m.modules = append(m.modules, v)
		}
	}

	return m
}

func (m *module) AddDecorator(decorators ...contracts.Decorator) contracts.Module {
	m.decorators = append(m.decorators, decorators...)
	return m
}

func (m *module) AddHook(hooks ...contracts.Hook) contracts.Module {
	m.hooks = append(m.hooks, hooks...)
	return m
}

func (m *module) AddInvoker(invokers ...contracts.Invoker) contracts.Module {
	m.invokers = append(m.invokers, invokers...)
	return m
}

func (m *module) AddModule(modules ...contracts.Module) contracts.Module {
	m.modules = append(m.modules, modules...)
	return m
}

func (m *module) AddProvider(providers ...contracts.Provider) contracts.Module {
	m.providers = append(m.providers, providers...)
	return m
}

func (m *module) Apply(container *dig.Container) error {
	for _, p := range m.providers {
		var opts []dig.ProvideOption
		if p.Name != "" {
			opts = append(opts, dig.Name(p.Name))
		}
		if p.Group != "" {
			opts = append(opts, dig.Group(p.Group))
		}

		if err := container.Provide(p.Fn, opts...); err != nil {
			return fmt.Errorf("module %q: provide: %w", m.name, err)
		}
	}

	for _, d := range m.decorators {
		if err := container.Decorate(d.Fn); err != nil {
			return fmt.Errorf("module %q: decorate: %w", m.name, err)
		}
	}

	for _, i := range m.invokers {
		if err := container.Invoke(i.Fn); err != nil {
			return fmt.Errorf("module %q: invoke: %w", m.name, err)
		}
	}

	for _, mod := range m.modules {
		if err := mod.Apply(container); err != nil {
			return err
		}
	}

	return nil
}

func (m *module) ApplyHooks(container *dig.Container, waiter waiterContracts.Waiter) error {
	for _, hookFn := range m.hooks {
		result, err := hookFn(container)
		if err != nil {
			return fmt.Errorf("module %q: hook error: %w", m.name, err)
		}

		if result.Name == "" {
			return fmt.Errorf("module %q: hook returned empty name", m.name)
		}

		hook := waiterContracts.WaiterHook{
			Name:    result.Name,
			Wait:    result.Wait,
			Cleanup: result.Cleanup,
		}
		if err := waiter.AddHook(hook); err != nil {
			return fmt.Errorf("module %q: hook: %w", m.name, err)
		}
	}

	for _, mod := range m.modules {
		if err := mod.ApplyHooks(container, waiter); err != nil {
			return err
		}
	}

	return nil
}

func (m *module) IsDependency() {}
