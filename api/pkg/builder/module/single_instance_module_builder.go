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

type SingleInstanceModuleOptions[Config configcontracts.Configurable, Instance, Dependencies any] struct {
	Name             string
	ConfigKey        string
	ConfigDecodeHook []mapstructure.DecodeHookFunc
	ProviderFn       func(context.Context, Config, Dependencies) (Instance, error)
	HookFn           func(context.Context, Instance) containercontracts.HookResolveResult
}

type SingleInstanceModuleResponse = func(
	context.Context,
	environmentenum.Environment,
	configcontracts.ConfigLoader,
) containercontracts.Module

func NewSingleInstanceModule[Config configcontracts.Configurable, Instance, Dependencies any](
	opts SingleInstanceModuleOptions[Config, Instance, Dependencies],
) SingleInstanceModuleResponse {
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

		var cfg Config
		if err := configLoader.LoadKey(env, &cfg, opts.ConfigKey); err != nil {
			log.Fatalf("Failed to build %q module with error: %v", opts.Name, err)
		}

		cfg.SetDefaults()
		if err := cfg.Validate(); err != nil {
			log.Fatalf("Failed to validate %q config with error: %v", opts.Name, err)
		}

		mod := container.NewModule(opts.Name, container.Provider(func() Config { return cfg }))

		mod.AddProvider(container.Provider(func(deps Dependencies) (Instance, error) {
			return opts.ProviderFn(ctx, cfg, deps)
		}))

		if opts.HookFn != nil {
			mod.AddHook(func(c *dig.Container) (containercontracts.HookResolveResult, error) {
				var instance Instance
				if err := c.Invoke(func(i Instance) { instance = i }); err != nil {
					return containercontracts.HookResolveResult{}, fmt.Errorf("%s hook: %w", opts.Name, err)
				}
				return opts.HookFn(ctx, instance), nil
			})
		}

		return mod
	}
}
