package httpserver

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/config"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/internal/adapters/echo"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/types/httpservertype"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func CreateHTTPServer(
	name string,
	cfg *config.HTTPServerOptions,
	logger loggercontracts.Logger,
) (contracts.HTTPServer, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	var adapter contracts.HTTPServer
	var err error

	switch cfg.Type {
	case httpservertype.HttpServerTypes.ECHO:
		adapter = echo.NewEchoAdapter(name, cfg, logger)
	default:
		err = fmt.Errorf("unsupported HTTP server type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	return adapter, nil
}
