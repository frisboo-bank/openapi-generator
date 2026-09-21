package otel

import (
	"context"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

var _ contracts.TracerAdapter = (*otelTracerAdapter)(nil)

type otelTracerAdapter struct {
	name           string
	logger         loggercontracts.Logger
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
}

func NewOtelTracerAdapter(name string, cfg *models.TracerOptions, resource *sdkresource.Resource, logger loggercontracts.Logger) (contracts.TracerAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("resource", resource)
	validation.AssertNotNil("logger", logger)

	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &otelTracerAdapter{
		name:           name,
		tracerProvider: tp,
		tracer:         tp.Tracer("openapi-generator-service"),
		logger:         logger,
	}, nil
}

func (o *otelTracerAdapter) Start(ctx context.Context, event string) (context.Context, contracts.TracerSpan) {
	ctx, span := o.tracer.Start(ctx, event)
	return ctx, &otelTracerSpan{span: span}
}

// Close implements [contracts.TracerAdapter].
func (o *otelTracerAdapter) Close(ctx context.Context) error {
	return o.tracerProvider.Shutdown(ctx)
}

func (o *otelTracerAdapter) Type() tracertype.TracerType {
	return tracertype.TracerTypes.OPEN_TELEMETRY
}
func (o *otelTracerAdapter) Name() string                   { return o.name }
func (o *otelTracerAdapter) Logger() loggercontracts.Logger { return o.logger }

// otelTracerSpan adapts [trace.Span] to [contracts.TracerSpan].
type otelTracerSpan struct {
	span trace.Span
}

func (s *otelTracerSpan) End()                  { s.span.End() }
func (s *otelTracerSpan) RecordError(err error) { s.span.RecordError(err) }
