package waiter

import (
	"context"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	"frisboo-bank/openapi-generator-service/pkg/waiter/config"
	"frisboo-bank/openapi-generator-service/pkg/waiter/contracts"
	waiterinternal "frisboo-bank/openapi-generator-service/pkg/waiter/internal"

	"go.uber.org/dig"
)

type WaiterModuleDependencies struct {
	dig.In
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var WaiterModule = module.NewSingleInstanceModule(
	module.SingleInstanceModuleOptions[*config.WaiterOptions, contracts.Waiter, WaiterModuleDependencies]{
		Name:      "waiter",
		ConfigKey: "waiter",
		ProviderFn: func(
			ctx context.Context,
			cfg *config.WaiterOptions,
			dependencies WaiterModuleDependencies,
		) (contracts.Waiter, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return waiterinternal.NewWaiter(ctx, cfg, loggerInstance)
		},
	},
)
