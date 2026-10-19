package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containercontracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/config"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/types"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/waiter/types/hookstatus"

	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	httpserverinternal "frisboo-bank/openapi-generator-service/pkg/http/http_server/internal"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type HTTPServerDependencies struct {
	dig.In
	Env     environmentenum.Environment
	Loggers module.DependenciesMap[loggercontracts.Logger]
}

var HTTPServerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.HTTPServerOptions, contracts.HTTPServer, HTTPServerDependencies]{
		Name:             "http-server",
		ConfigKey:        "http-servers",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{types.HTTPServerEnumsDecodeHook()},
		ProviderFn: func(name string, cfg *config.HTTPServerOptions, dependencies HTTPServerDependencies) (contracts.HTTPServer, error) {
			loggerInstance, err := dependencies.Loggers.Get(cfg.Logger)
			if err != nil {
				return nil, err
			}

			return httpserverinternal.CreateHTTPServer(name, cfg, loggerInstance, dependencies.Env)
		},
		HookFn: func(name string, instance contracts.HTTPServer) containercontracts.HookResolveResult {
			return containercontracts.HookResolveResult{
				Name: "http-server:" + name,
				Start: func(ctx context.Context) error {
					errCh := make(chan error, 1)
					go func() { errCh <- instance.Start(ctx) }()
					select {
					case err := <-errCh:
						if err != nil && !errors.Is(err, http.ErrServerClosed) {
							return fmt.Errorf("http-server %q stopped unexpectedly: %w", name, err)
						}
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				},
				Stop: func(ctx context.Context) error {
					if err := instance.Stop(ctx); err != nil {
						if !errors.Is(err, http.ErrServerClosed) {
							instance.Logger().Errorf("http-server %q shutdown failed: %v", name, err)
							return err
						}
					}
					instance.Logger().Infof("http-server %q shut down gracefully", name)
					return nil
				},
				Status: func() hookstatus.HookStatus {
					if instance == nil {
						return hookstatus.HookStatuses.UNKNOWN
					}
					return instance.Status()
				},
			}
		},
	},
)
