package contracts

import "net/http"

type RouterEngine interface {
	Group(path string, middlewares ...any) RouterEngine
	Handle(method string, path string, handler any, middlewares ...any)
	Static(prefix string, root string)
}

type RouteBuilder interface {
	Root() RouteGroup
}

type RouteGroup interface {
	// Nested group
	Group(path string, middlewares ...any) RouteGroup

	DELETE(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	GET(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	HEAD(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	OPTIONS(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	PATCH(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	POST(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	PUT(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	TRACE(path string, handler http.HandlerFunc, middlewares ...any) RouteGroup

	Handle(method string, path string, handler http.HandlerFunc, middlewares ...any) RouteGroup
	Static(prefix string, root string) RouteGroup
}
