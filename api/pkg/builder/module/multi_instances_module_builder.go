package module

import (
	"context"
	"fmt"
	"log"

	configcontracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type MultiInstancesModuleOptions[Config configcontracts.Configurable, Instance, Dependencies any] struct {
	Name             string
	ConfigKey        string
	ConfigDecodeHook []mapstructure.DecodeHookFunc
	ProviderFn       func(context.Context, string, Config, Dependencies) (Instance, error)
	HookFn           func(context.Context, string, Instance) containercontracts.HookResolveResult
}

type MultiInstancesModuleResponse = func(
	context.Context,
	environmentenum.Environment,
	configcontracts.ConfigLoader,
) containercontracts.Module

func NewMultiInstancesModule[Config configcontracts.Configurable, Instance, Dependencies any](
	opts MultiInstancesModuleOptions[Config, Instance, Dependencies],
) MultiInstancesModuleResponse {
	validation.AssertNotEmpty("name", opts.Name)
	validation.AssertNotEmpty("configKey", opts.ConfigKey)
	validation.AssertNotNil("ProviderFn", opts.ProviderFn)

	return func(
		ctx context.Context,
		env environmentenum.Environment,
		configLoader configcontracts.ConfigLoader,
	) containercontracts.Module {
		validation.AssertNotNil("ctx", ctx)
		validation.AssertNotNil("env", env)
		validation.AssertNotNil("configLoader", configLoader)

		if len(opts.ConfigDecodeHook) > 0 {
			configLoader.RegisterDecodeHookFunc(opts.ConfigDecodeHook...)
		}

		type ConfigsMapType = map[string]Config
		type InstancesMapType = DependenciesMap[Instance]

		var cfgMap ConfigsMapType
		if err := configLoader.LoadKey(env, &cfgMap, opts.ConfigKey); err != nil {
			log.Fatalf("Failed to build %q module with error: %v", opts.Name, err)
		}

		mod := container.NewModule(opts.Name, container.Provider(func() ConfigsMapType { return cfgMap }))
		if len(cfgMap) == 0 {
			return mod
		}

		mod.AddProvider(container.Provider(
			func(dependencies Dependencies) (InstancesMapType, error) {
				instances := make(InstancesMapType)

				for name, cfg := range cfgMap {
					if !cfg.GetEnabled() {
						continue
					}
					cfg.SetDefaults()
					if err := cfg.Validate(); err != nil {
						return nil, fmt.Errorf("validate config %q: %w", name, err)
					}

					instance, err := opts.ProviderFn(ctx, name, cfg, dependencies)
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

		for name := range cfgMap {
			mod.AddProvider(containercontracts.Provider{
				Fn:   func(instances InstancesMapType) Instance { return instances[name] },
				Name: opts.Name + ":" + name,
			})

			if opts.HookFn != nil {
				mod.AddHook(func(c *dig.Container) (containercontracts.HookResolveResult, error) {
					var all InstancesMapType
					if err := c.Invoke(func(m InstancesMapType) { all = m }); err != nil {
						return containercontracts.HookResolveResult{}, fmt.Errorf("%s %q hook: %w", opts.Name, name, err)
					}
					instance, ok := all[name]
					if !ok {
						return containercontracts.HookResolveResult{}, fmt.Errorf("%s %q hook: instance not found", opts.Name, name)
					}
					return opts.HookFn(ctx, name, instance), nil
				})
			}
		}

		return mod
	}
}
