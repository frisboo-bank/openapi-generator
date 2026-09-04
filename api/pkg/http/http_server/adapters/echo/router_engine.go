package echo

import (
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	echoVendor "github.com/labstack/echo/v4"
)

var _ contracts.RouterEngine = (*echoRouterEngine)(nil)

type echoRouterEngine struct {
	echo   *echoVendor.Echo
	group  *echoVendor.Group
	logger loggerContracts.Logger
	prefix string
}

func newRouterEngine(e *echoVendor.Echo, logger loggerContracts.Logger) contracts.RouterEngine {
	validation.AssertNotNil("echo", e)
	validation.AssertNotNil("logger", logger)

	return &echoRouterEngine{
		echo:   e,
		group:  nil,
		logger: logger,
	}
}

func (e *echoRouterEngine) Group(path string, middlewares ...any) contracts.RouterEngine {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrapf(err, "invalid middleware for group:{%s}", path))
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	var g *echoVendor.Group
	if e.group != nil {
		g = e.group.Group(path, ms...)
	} else {
		g = e.echo.Group(path, ms...)
	}

	newEngine := &echoRouterEngine{
		echo:   e.echo,
		group:  g,
		logger: e.logger,
		prefix: e.getFullPath(path),
	}

	e.logger.Debugf("route group registered prefix:{%s} middlewares:{%d}", newEngine.prefix, len(ms))

	return newEngine
}

func (e *echoRouterEngine) Handle(method string, path string, handler any, middlewares ...any) {
	validation.AssertNotEmpty("method", method)
	validation.AssertNotNil("handler", handler)

	method = strings.ToUpper(method)

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	h, err := ToHandlerFunc(handler)
	if err != nil {
		panic(syserrors.Wrapf(err, "invalid handler for route method:{%s} path:{%s}", method, path))
	}

	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrapf(err, "invalid middleware for route method:{%s} path:{%s}", method, path))
	}

	if e.group != nil {
		e.group.Add(method, path, h, ms...)
	} else {
		e.echo.Add(method, path, h, ms...)
	}

	e.logger.Debugf("route registered method:{%s} path:{%s} middlewares:{%d}", method, e.getFullPath(path), len(ms))
}

func (e *echoRouterEngine) Static(prefix string, root string) {
	validation.AssertNotEmpty("root", root)

	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	if e.group != nil {
		e.group.Static(prefix, root)
	} else {
		e.echo.Static(prefix, root)
	}

	e.logger.Debugf("static route registered prefix:{%s} root:{%s} path:{%s}", prefix, root, e.getFullPath(prefix))
}

func (e *echoRouterEngine) getFullPath(path string) string {
	prefix := strings.TrimRight(e.prefix, "/")
	path = strings.TrimLeft(path, "/")

	if e.prefix == "" {
		return "/" + path
	}
	return prefix + "/" + path
}
