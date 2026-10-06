package noop

import (
	"time"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metrictype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"
)

var _ contracts.MetricsAdapter = (*noopMetricsAdapter)(nil)

type noopMetricsAdapter struct {
	name   string
	logger loggercontracts.Logger
}

func NewNoopMetricsAdapter(
	name string,
	logger loggercontracts.Logger,
) (contracts.MetricsAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("logger", logger)

	return &noopMetricsAdapter{
		name:   name,
		logger: logger,
	}, nil
}

func (o *noopMetricsAdapter) RecordDuration(name string, duraction time.Duration, attrs ...any) {}

func (o *noopMetricsAdapter) Close() error { return nil }

func (o *noopMetricsAdapter) Logger() loggercontracts.Logger { return o.logger }
func (o *noopMetricsAdapter) Name() string                   { return o.name }
func (o *noopMetricsAdapter) Type() metrictype.MetricsType   { return metrictype.MetricsTypes.NOOP }
