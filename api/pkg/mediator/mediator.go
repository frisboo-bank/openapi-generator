package mediator

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator/models"

	"github.com/hashicorp/go-multierror"
)

var _ contracts.Mediator = (*mediator)(nil)

type mediator struct {
	mu            sync.RWMutex
	handlers      map[reflect.Type]func(context.Context, any) (any, error)
	notifications map[reflect.Type][]func(context.Context, any) error
}

func NewMediator(
	cfg *models.MediatorOptions,
	logger loggerContracts.Logger,
) (contracts.Mediator, error) {
	return &mediator{
		handlers:      map[reflect.Type]func(context.Context, any) (any, error){},
		notifications: map[reflect.Type][]func(context.Context, any) error{},
	}, nil
}

func (m *mediator) Publish(ctx context.Context, notification any) error {
	t := reflect.TypeOf(notification)

	m.mu.RLock()
	fns, ok := m.notifications[t]
	m.mu.RUnlock()

	if !ok || len(fns) == 0 {
		return nil
	}

	var errs error
	for _, fn := range fns {
		errs = multierror.Append(errs, fn(ctx, notification))
	}

	return errs
}

func (m *mediator) Send(ctx context.Context, request any) (any, error) {
	t := reflect.TypeOf(request)

	m.mu.RLock()
	fn, ok := m.handlers[t]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("mediator: no handler registered for %T", request)
	}

	return fn(ctx, request)
}

func (m *mediator) registerHandler(
	requestType reflect.Type,
	fn func(context.Context, any) (any, error),
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.handlers[requestType]; exists {
		return fmt.Errorf("mediator: handler already registered for %v", requestType)
	}

	m.handlers[requestType] = fn

	return nil
}
