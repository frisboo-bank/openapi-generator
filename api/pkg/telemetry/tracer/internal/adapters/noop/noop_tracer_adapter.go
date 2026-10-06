package noop

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
)

var (
	_ contracts.TracerAdapter = (*noopTracerAdapter)(nil)
	_ contracts.TracerSpan    = (*noopTracerSpan)(nil)
)

type noopTracerSpan struct{}

func (n *noopTracerSpan) RecordError(err error) {}
func (n *noopTracerSpan) End()                  {}

type noopTracerAdapter struct {
	name   string
	logger loggercontracts.Logger
}

func NewNoopTracerAdapter(
	name string,
	logger loggercontracts.Logger,
) (contracts.TracerAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("logger", logger)

	return &noopTracerAdapter{
		name:   name,
		logger: logger,
	}, nil
}

func (n *noopTracerAdapter) Start(ctx context.Context, event string) (context.Context, contracts.TracerSpan) {
	return ctx, &noopTracerSpan{}
}

func (n *noopTracerAdapter) Close(ctx context.Context) error { return nil }

func (n *noopTracerAdapter) Logger() loggercontracts.Logger { return n.logger }

func (n *noopTracerAdapter) Name() string                { return n.name }
func (n *noopTracerAdapter) Type() tracertype.TracerType { return tracertype.TracerTypes.NOOP }
