package mediator

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator/models"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ contracts.Mediator = (*mediator)(nil)

type mediator struct {
	mu                   sync.RWMutex
	requestHandlers      map[reflect.Type]any
	notificationHandlers map[reflect.Type][]any
	logger               loggercontracts.Logger
}

func NewMediator(
	cfg *models.MediatorOptions,
	logger loggercontracts.Logger,
) (contracts.Mediator, error) {
	return &mediator{
		logger:               logger,
		requestHandlers:      make(map[reflect.Type]any, 0),
		notificationHandlers: make(map[reflect.Type][]any, 0),
	}, nil
}

func (m *mediator) RegisterRequest(requestType contracts.TRequest, handler any) error {
	validation.AssertNotNil("handler", handler)

	t := reflect.TypeOf(requestType)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.requestHandlers[t]; exists {
		return fmt.Errorf("mediator: handler already registered for request type %v", t)
	}
	m.requestHandlers[t] = handler

	return nil
}

func (m *mediator) RegisterNotification(notificationType contracts.TNotification, handler any) error {
	validation.AssertNotNil("handler", handler)

	t := reflect.TypeOf(notificationType)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.notificationHandlers[t]; !exists {
		m.notificationHandlers[t] = make([]any, 0)
	}
	m.notificationHandlers[t] = append(m.notificationHandlers[t], handler)

	return nil
}

func (m *mediator) Send(ctx context.Context, request contracts.Request) (contracts.Response, error) {
	validation.AssertNotNil("request", request)

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("mediator: context cancelled before send: %w", err)
	}

	t := reflect.TypeOf(request)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	m.mu.RLock()
	handlerRaw, exists := m.requestHandlers[t]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("mediator: no handler registered for %T", request)
	}

	handler, ok := handlerRaw.(func(context.Context, any) (any, error))
	if !ok {
		return nil, fmt.Errorf("mediator: invalid handler type stored for %T", request)
	}

	return handler(ctx, request)
}

func (m *mediator) Publish(ctx context.Context, notification contracts.Notification) error {
	panic("unimplemented")
}
