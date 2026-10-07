package contracts

import (
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
)

type (
	Attributes map[string]string

	Metrics interface {
		MetricsAdapter
	}

	MetricsAdapter interface {
		RecordDuration(name string, duraction time.Duration, attrs ...any)
		Close() error
		Name() string
		Type() metricstype.MetricsType
		Logger() loggercontracts.Logger
	}
)
