package echo

import (
	"context"
	"net"
	"strings"
	"sync/atomic"
	"time"

	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/config"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/internal/adapters/echo/middlewares/correlationid"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/internal/adapters/echo/middlewares/logger"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/routing"
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/types/httpservertype"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"frisboo-bank/openapi-generator-service/pkg/waiter/types/hookstatus"

	echovendor "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ contracts.HTTPServer = (*echoAdapter)(nil)

type echoAdapter struct {
	address       string
	bodyLimit     string
	gzipLevel     int
	ignoreLogUrls []string
	logger        loggerContracts.Logger
	name          string
	routeBuilder  contracts.RouteBuilder
	server        *echovendor.Echo

	isStarted atomic.Bool
}

func NewEchoAdapter(
	name string,
	cfg *config.HTTPServerOptions,
	logger loggerContracts.Logger,
	env environmentenum.Environment,
) contracts.HTTPServer {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)
	validation.AssertValidEnum("env", env)

	srv := echovendor.New()
	srv.HideBanner = true

	srv.Server.ReadTimeout = cfg.ReadTimeout
	srv.Server.ReadHeaderTimeout = cfg.ReadHeaderTimeout
	srv.Server.WriteTimeout = cfg.WriteTimeout
	srv.Server.IdleTimeout = cfg.IdleTimeout
	srv.Server.MaxHeaderBytes = cfg.MaxHeaderBytes

	routerEngine := newRouterEngine(srv, logger)
	routeBuilder := routing.NewRouteBuilder(routerEngine)

	return &echoAdapter{
		address:       cfg.Address(),
		bodyLimit:     cfg.BodyLimit,
		gzipLevel:     cfg.GzipLevel,
		ignoreLogUrls: cfg.IgnoreLogUrls,
		logger:        logger,
		name:          name,
		routeBuilder:  routeBuilder,
		server:        srv,
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
		correlationid.NewEchoLoggerMiddleware(e.logger),
		middleware.Recover(),
		middleware.RemoveTrailingSlash(),
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
	e.isStarted.Store(true)
	return e.server.Start(e.address)
}

func (e *echoAdapter) Status() hookstatus.HookStatus {
	switch {
	case !e.isStarted.Load():
		return hookstatus.HookStatuses.IDLE
	case e.isListening():
		return hookstatus.HookStatuses.RUNNING
	default:
		return hookstatus.HookStatuses.STOPPED
	}
}

func (e *echoAdapter) Stop(ctx context.Context) error {
	return e.server.Shutdown(ctx)
}

func (e *echoAdapter) isListening() bool {
	network := "tcp"
	addr := e.address

	if strings.HasPrefix(addr, "/") || strings.HasPrefix(addr, "unix:") {
		network = "unix"
		addr = strings.TrimPrefix(addr, "unix:")
	}

	conn, err := net.DialTimeout(network, addr, 200*time.Millisecond)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func (e *echoAdapter) ListRoutes() []any {
	rs := e.server.Routes()
	out := make([]any, len(rs))
	for i, r := range rs {
		out[i] = r
	}
	return out
}

func (e *echoAdapter) skipper() func(c echovendor.Context) bool {
	return func(c echovendor.Context) bool {
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
