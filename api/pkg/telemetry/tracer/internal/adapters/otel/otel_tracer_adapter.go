package otel

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types/tracertype"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

var _ contracts.Tracer = (*otelTracerAdapter)(nil)

type otelTracerAdapter struct {
	name           string
	logger         loggercontracts.Logger
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
}

func NewOtelTracerAdapter(
	ctx context.Context,
	name string,
	cfg *config.TracerOptions,
	resource *sdkresource.Resource,
	logger loggercontracts.Logger,
) (contracts.Tracer, error) {
	validation.AssertNotNil("ctx", ctx)
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("ctx", ctx)
	validation.AssertNotNil("resource", resource)
	validation.AssertNotNil("logger", logger)

	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.Endpoint)}
	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
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

func (s *otelTracerSpan) End() { s.span.End() }

func (s *otelTracerSpan) RecordError(err error) { s.span.RecordError(err) }

func (s *otelTracerSpan) RecordPanic(v any) {
	s.span.RecordError(fmt.Errorf("panic: %v", v))
}

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
