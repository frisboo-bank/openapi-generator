package routing

import (
	"frisboo-bank/openapi-generator-service/pkg/http/http_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ contracts.RouteBuilder = (*routeBuilder)(nil)

type routeBuilder struct {
	engine contracts.RouterEngine
}

func NewRouteBuilder(engine contracts.RouterEngine) contracts.RouteBuilder {
	validation.AssertNotNil("engine", engine)

	return &routeBuilder{engine: engine}
}

func (r *routeBuilder) Root() contracts.RouteGroup {
	return newRouteGroup(r.engine)
}
