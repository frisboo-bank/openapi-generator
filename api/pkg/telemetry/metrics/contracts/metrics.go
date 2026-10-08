package contracts

import (
	"context"
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
)

type (
	Attributes map[string]string

	Metrics interface {
		MetricsAdapter
	}

	MetricsAdapter interface {
		RecordDuration(name string, duraction time.Duration, attrs ...any)
		Close(ctx context.Context) error
		Name() string
		Type() metricstype.MetricsType
		Logger() loggercontracts.Logger
	}
)
