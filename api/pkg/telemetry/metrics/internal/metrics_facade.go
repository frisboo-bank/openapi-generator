package metrics

import (
	"time"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
)

var _ contracts.Metrics = (*metrics)(nil)

type metrics struct {
	adapter contracts.MetricsAdapter
}

func (m *metrics) RecordDuration(name string, duraction time.Duration, attrs ...any) {
	m.adapter.RecordDuration(name, duraction, attrs...)
}

func (m *metrics) Close() error { return m.adapter.Close() }

func (m *metrics) Name() string                   { return m.adapter.Name() }
func (m *metrics) Type() metricstype.MetricsType  { return m.adapter.Type() }
func (m *metrics) Logger() loggercontracts.Logger { return m.adapter.Logger() }
