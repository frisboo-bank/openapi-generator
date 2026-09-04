package module

import (
	"fmt"
	"log"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"go.uber.org/dig"
)

type MultiInstancesModuleOptions[
	Config configContracts.Configurable,
	Instance any,
	Extra any,
] struct {
	Name       string
	ConfigKey  string
	ProviderFn func(name string, cfg Config, env environmentEnum.Environment, logger loggerContracts.Logger, extra Extra) (Instance, error)
	HookFn     func(name string, instance Instance) containerContracts.HookResolveResult
}

type MultiInstancesModuleResponse = func(
	env environmentEnum.Environment,
	configLoader configContracts.ConfigLoader,
) containerContracts.Module

func NewMultiInstancesModule[Config configContracts.Configurable, Instance any, Extra any](
	opts MultiInstancesModuleOptions[Config, Instance, Extra],
) MultiInstancesModuleResponse {
	validation.AssertNotEmpty("name", opts.Name)
	validation.AssertNotEmpty("configKey", opts.ConfigKey)
	validation.AssertNotNil("ProviderFn", opts.ProviderFn)

	return func(
		env environmentEnum.Environment,
		configLoader configContracts.ConfigLoader,
	) containerContracts.Module {
		validation.AssertNotNil("env", env)
		validation.AssertNotNil("configLoader", configLoader)

		type ConfigsMapType = map[string]Config
		type InstancesMapType = map[string]Instance

		var cfgMap ConfigsMapType
		if err := configLoader.LoadKey(env, &cfgMap, opts.ConfigKey); err != nil {
			log.Fatalf("Failed to build %q module with error: %v", opts.Name, err)
		}

		mod := container.NewModule(
			opts.Name,
			container.Provider(func() ConfigsMapType { return cfgMap }),
		)

		if len(cfgMap) == 0 {
			return mod
		}

		mod.AddProvider(container.Provider(
			func(loggers map[string]loggerContracts.Logger, extra Extra) (InstancesMapType, error) {
				instances := make(InstancesMapType)

				for name, cfg := range cfgMap {
					if !cfg.GetEnabled() {
						continue
					}

					cfg.SetDefaults()
					if err := cfg.Validate(); err != nil {
						return nil, fmt.Errorf("validate config %q: %w", name, err)
					}

					loggerName := cfg.GetLogger()
					if loggerName == "" {
						loggerName = name
					}
					logger, ok := loggers[loggerName]
					if !ok {
						return nil, fmt.Errorf("logger %q not found for %q %q", loggerName, opts.Name, name)
					}

					instance, err := opts.ProviderFn(name, cfg, env, logger, extra)
					if err != nil {
						return nil, fmt.Errorf("%s %q: %w", opts.Name, name, err)
					}

					if any(instance) == nil {
						continue
					}

					instances[name] = instance
				}

				return instances, nil
			},
		))

		for name, cfg := range cfgMap {
			if !cfg.GetEnabled() {
				continue
			}

			mod.AddProvider(containerContracts.Provider{
				Fn:   func(instances InstancesMapType) Instance { return instances[name] },
				Name: opts.Name + ":" + name,
			})

			if opts.HookFn != nil {
				n := name
				mod.AddHook(func(c *dig.Container) (containerContracts.HookResolveResult, error) {
					var all InstancesMapType
					if err := c.Invoke(func(m InstancesMapType) { all = m }); err != nil {
						return containerContracts.HookResolveResult{},
							fmt.Errorf("%s %q hook: %w", opts.Name, n, err)
					}
					instance, ok := all[n]
					if !ok {
						return containerContracts.HookResolveResult{},
							fmt.Errorf("%s %q hook: instance not found", opts.Name, n)
					}
					return opts.HookFn(n, instance), nil
				})
			}
		}

		return mod
	}
}
