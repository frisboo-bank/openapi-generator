package contracts

import (
	"context"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
)

type (
	Tracer interface {
		TracerAdapter
	}

	TracerAdapter interface {
		Start(ctx context.Context, event string) (context.Context, TracerSpan)
		Close(ctx context.Context) error
		Name() string
		Type() tracertype.TracerType
		Logger() loggercontracts.Logger
	}

	TracerSpan interface {
		End()
		RecordError(err error)
	}
)
