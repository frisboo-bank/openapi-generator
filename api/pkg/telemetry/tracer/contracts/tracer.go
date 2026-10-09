package contracts

import (
	"context"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types/tracertype"
)

type (
	Tracer interface {
		Start(ctx context.Context, event string) (context.Context, TracerSpan)
		Close(ctx context.Context) error
		Name() string
		Type() tracertype.TracerType
		Logger() loggercontracts.Logger
	}

	TracerSpan interface {
		End()
		RecordError(err error)
		RecordPanic(v any)
	}
)
