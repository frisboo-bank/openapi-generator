package otel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models"
	metrictype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	metric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

var _ contracts.MetricsAdapter = (*otelMetricsAdapter)(nil)

type otelMetricsAdapter struct {
	name             string
	meterProvider    *sdkmetric.MeterProvider
	meter            metric.Meter
	histogram        sync.Once
	float64Histogram metric.Float64Histogram
	logger           loggercontracts.Logger
}

func NewOtelMetricsAdapter(name string, cfg *models.MetricsOptions, resource *sdkresource.Resource, logger loggercontracts.Logger) (contracts.MetricsAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	ctx := context.Background()

	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(resource),
	)

	otel.SetMeterProvider(mp)
	return &otelMetricsAdapter{
		name:          name,
		meterProvider: mp,
		meter:         mp.Meter("openapi-generator-service"),
		logger:        logger,
	}, nil
}

func (o *otelMetricsAdapter) RecordDuration(name string, duraction time.Duration, attrs ...any) {
	attrs = append(attrs, "duration", duraction)
	measurement := metric.WithAttributeSet(toAttributeSet(attrs))

	o.histogram.Do(func() {
		hist, err := o.meter.Float64Histogram(name)
		if err != nil {
			o.logger.Errorf("failed to create histogram %q: %v", name, err)
			return
		}
		o.float64Histogram = hist
	})

	if o.float64Histogram != nil {
		o.float64Histogram.Record(context.Background(), float64(duraction), measurement)
	}
}

func toAttributeSet(attrs []any) attribute.Set {
	var set []attribute.KeyValue
	for i := 0; i+1 < len(attrs); i += 2 {
		key, ok := attrs[i].(string)
		if !ok {
			continue
		}
		set = append(set, attribute.String(key, fmt.Sprint(attrs[i+1])))
	}
	return attribute.NewSet(set...)
}

func (o *otelMetricsAdapter) Close(ctx context.Context) error {
	return o.meterProvider.Shutdown(ctx)
}

func (o *otelMetricsAdapter) Type() metrictype.MetricsType {
	return metrictype.MetricsTypes.OPEN_TELEMETRY
}
func (o *otelMetricsAdapter) Name() string                   { return o.name }
func (o *otelMetricsAdapter) Logger() loggercontracts.Logger { return o.logger }
