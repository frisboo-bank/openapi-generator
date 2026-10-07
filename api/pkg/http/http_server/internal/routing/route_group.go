package routing

import (
	"net/http"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ contracts.RouteGroup = (*routeGroup)(nil)

type routeGroup struct {
	engine contracts.RouterEngine
}

func newRouteGroup(engine contracts.RouterEngine) contracts.RouteGroup {
	validation.AssertNotNil("engine", engine)

	return &routeGroup{engine: engine}
}

func (r *routeGroup) Group(path string, middlewares ...any) contracts.RouteGroup {
	g := r.engine.Group(path, middlewares...)
	return newRouteGroup(g)
}

func (r *routeGroup) DELETE(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodDelete, path, handler, middlewares...)
}

func (r *routeGroup) GET(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodGet, path, handler, middlewares...)
}

func (r *routeGroup) HEAD(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodHead, path, handler, middlewares...)
}

func (r *routeGroup) OPTIONS(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodOptions, path, handler, middlewares...)
}

func (r *routeGroup) PATCH(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodPatch, path, handler, middlewares...)
}

func (r *routeGroup) POST(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodPost, path, handler, middlewares...)
}

func (r *routeGroup) PUT(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodPut, path, handler, middlewares...)
}

func (r *routeGroup) TRACE(path string, handler http.HandlerFunc, middlewares ...any) contracts.RouteGroup {
	return r.Handle(http.MethodTrace, path, handler, middlewares...)
}

func (r *routeGroup) Handle(
	method string,
	path string,
	handler http.HandlerFunc,
	middlewares ...any,
) contracts.RouteGroup {
	r.engine.Handle(method, path, handler, middlewares...)
	return r
}

func (r *routeGroup) Static(prefix string, root string) contracts.RouteGroup {
	r.engine.Static(prefix, root)
	return r
}
