package httpserver

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/internal/adapters/echo"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/models"
	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func CreateHTTPServer(
	name string,
	cfg *models.HTTPServerOptions,
	logger loggercontracts.Logger,
	env environmentEnum.Environment,
) (contracts.HTTPServer, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)
	validation.AssertValidEnum("env", env)

	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	switch cfg.Type {
	case httpservertype.HttpServerTypes.ECHO:
		return &httpServer{echo.NewEchoAdapter(name, cfg, logger, env)}, nil
	default:
		return nil, syserrors.Newf("no http server of type %q exists", cfg.Type)
	}
}
