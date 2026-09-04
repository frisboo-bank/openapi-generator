package echo

import (
	"context"
	"strings"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/adapters/echo/middlewares/logger"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/models"
	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/routing"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"

	echoVendor "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ contracts.HTTPServer = (*echoAdapter)(nil)

type echoAdapter struct {
	address       string
	bodyLimit     string
	server        *echoVendor.Echo
	logger        loggerContracts.Logger
	name          string
	routeBuilder  contracts.RouteBuilder
	ignoreLogUrls []string
	gzipLevel     int
}

func NewEchoAdapter(
	name string,
	cfg *models.HTTPServerOptions,
	logger loggerContracts.Logger,
	env environmentEnum.Environment,
) contracts.HTTPServer {
	srv := echoVendor.New()
	srv.HideBanner = true

	srv.Server.ReadTimeout = cfg.ReadTimeout
	srv.Server.ReadHeaderTimeout = cfg.ReadHeaderTimeout
	srv.Server.WriteTimeout = cfg.WriteTimeout
	srv.Server.IdleTimeout = cfg.IdleTimeout
	srv.Server.MaxHeaderBytes = cfg.MaxHeaderBytes

	routerEngine := newRouterEngine(srv, logger)
	routeBuilder := routing.NewRouteBuilder(routerEngine)

	return &echoAdapter{
		name:          name,
		server:        srv,
		logger:        logger,
		routeBuilder:  routeBuilder,
		address:       cfg.Address(),
		bodyLimit:     cfg.BodyLimit,
		ignoreLogUrls: cfg.IgnoreLogUrls,
		gzipLevel:     cfg.GzipLevel,
	}
}

func (e *echoAdapter) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrap(err, "invalid middleware"))
	}
	e.server.Use(ms...)
}

func (e *echoAdapter) SetupDefaultMiddlewares() {
	skipper := e.skipper()

	e.server.Use(
		logger.NewEchoLoggerMiddleware(e.logger, logger.EchoLoggerMiddlewareOptions{
			Skipper: skipper,
		}),
		middleware.Recover(),
		middleware.AddTrailingSlash(),
		middleware.BodyLimit(e.bodyLimit),
		middleware.RequestID(),
		middleware.RequestLogger(),
		middleware.Secure(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: skipper,
			Level:   e.gzipLevel,
		}),
		middleware.CORS(),
	)
}

func (e *echoAdapter) Start(ctx context.Context) error {
	return e.server.Start(e.address)
}

func (e *echoAdapter) Stop(ctx context.Context) error {
	return e.server.Shutdown(ctx)
}

func (e *echoAdapter) ListRoutes() []any {
	rs := e.server.Routes()
	out := make([]any, len(rs))
	for i, r := range rs {
		out[i] = r
	}
	return out
}

func (e *echoAdapter) skipper() func(c echoVendor.Context) bool {
	return func(c echoVendor.Context) bool {
		rPath := c.Request().URL.Path
		for _, skip := range e.ignoreLogUrls {
			if strings.Contains(rPath, skip) {
				return true
			}
		}
		return false
	}
}

func (e *echoAdapter) Name() string { return e.name }
func (e *echoAdapter) Type() httpservertype.HttpServerType {
	return httpservertype.HttpServerTypes.ECHO
}
func (e *echoAdapter) RouteBuilder() contracts.RouteBuilder { return e.routeBuilder }
func (e *echoAdapter) Logger() loggerContracts.Logger       { return e.logger }
