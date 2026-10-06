package noop

import (
	"context"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metrictype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"
)

var _ contracts.MetricsAdapter = (*otelMetricsAdapter)(nil)

type otelMetricsAdapter struct {
	name   string
	ctx    context.Context
	logger loggercontracts.Logger
}

func NewNoopMetricsAdapter(
	name string,
	cfg *models.MetricsOptions,
	ctx context.Context,
	logger loggercontracts.Logger,
) (contracts.MetricsAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("ctx", ctx)
	validation.AssertNotNil("logger", logger)

	return &otelMetricsAdapter{
		name:   name,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (o *otelMetricsAdapter) RecordDuration(name string, duraction time.Duration, attrs ...any) {}

func (o *otelMetricsAdapter) Close() error { return nil }

func (o *otelMetricsAdapter) Logger() loggercontracts.Logger { return o.logger }
func (o *otelMetricsAdapter) Name() string                   { return o.name }
func (o *otelMetricsAdapter) Type() metrictype.MetricsType   { return metrictype.MetricsTypes.NOOP }
