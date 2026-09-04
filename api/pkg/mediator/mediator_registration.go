package mediator

import (
	"context"
	"fmt"
	"reflect"

	"frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
)

type HandlerRegistrar any

type HandlerRegistration struct {
	Register func(registrar HandlerRegistrar) error
}

func RegisterTypedHandler[TRequest any, TResponse any](
	registrar HandlerRegistrar,
	handler contracts.Handler[TRequest, TResponse],
) error {
	m, ok := registrar.(*mediator)
	if !ok {
		return fmt.Errorf("mediator: invalid registrar type")
	}

	var req TRequest
	requestType := reflect.TypeOf(req)

	fn := func(ctx context.Context, request any) (any, error) {
		return handler.Handle(ctx, request.(TRequest))
	}

	return m.registerHandler(requestType, fn)
}
