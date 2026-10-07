package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/openapi-generator-service/pkg/builder/module"
	containerContracts "frisboo-bank/openapi-generator-service/pkg/container/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/config"
	httpserverinternal "frisboo-bank/openapi-generator-service/pkg/http/http_server/internal"
	httpserverenums "frisboo-bank/openapi-generator-service/pkg/http/http_server/types/enums"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type HTTPServerDependencies struct {
	dig.In
}

var HTTPServerModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.HTTPServerOptions, contracts.HTTPServer, HTTPServerDependencies]{
		Name:      "http-server",
		ConfigKey: "http-servers",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{
			httpserverenums.HTTPServerEnumsDecodeHook(),
		},
		ProviderFn: func(
			name string,
			cfg *config.HTTPServerOptions,
			env environmentEnum.Environment,
			logger loggerContracts.Logger,
			extra HTTPServerDependencies,
		) (contracts.HTTPServer, error) {
			return httpserverinternal.CreateHTTPServer(name, cfg, logger, env)
		},
		HookFn: func(name string, instance contracts.HTTPServer) containerContracts.HookResolveResult {
			return containerContracts.HookResolveResult{
				Name: "http-server:" + name,
				Wait: func(ctx context.Context) error {
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
				Cleanup: func(ctx context.Context) error {
					if err := instance.Stop(ctx); err != nil {
						if !errors.Is(err, http.ErrServerClosed) {
							instance.Logger().Errorf("http-server %q shutdown failed: %v", name, err)
							return err
						}
					}

					instance.Logger().Infof("http-server %q shut down gracefully", name)

					return nil
				},
			}
		},
	},
)
