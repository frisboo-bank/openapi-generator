package tracer

import (
	"context"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
)

var _ contracts.Tracer = (*tracer)(nil)

type tracer struct {
	adapter contracts.TracerAdapter
}

func (t *tracer) Close(ctx context.Context) error {
	return t.adapter.Close(ctx)
}

func (t *tracer) Logger() loggercontracts.Logger {
	return t.adapter.Logger()
}

func (t *tracer) Name() string {
	return t.adapter.Name()
}

func (t *tracer) Start(ctx context.Context, event string) (context.Context, contracts.TracerSpan) {
	return t.adapter.Start(ctx, event)
}

func (t *tracer) Type() tracertype.TracerType {
	return t.adapter.Type()
}
