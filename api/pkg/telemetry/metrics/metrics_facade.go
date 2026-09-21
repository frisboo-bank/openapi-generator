package metrics

import (
	"context"
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"
)

var _ contracts.Metrics = (*metrics)(nil)

type metrics struct {
	adapter contracts.MetricsAdapter
}

func (m *metrics) RecordDuration(name string, duraction time.Duration, attrs ...any) {
	m.adapter.RecordDuration(name, duraction, attrs...)
}

func (m *metrics) Close(ctx context.Context) error { return m.adapter.Close(ctx) }
func (m *metrics) Name() string                    { return m.adapter.Name() }
func (m *metrics) Type() metricstype.MetricsType   { return m.adapter.Type() }
func (m *metrics) Logger() loggercontracts.Logger  { return m.adapter.Logger() }
