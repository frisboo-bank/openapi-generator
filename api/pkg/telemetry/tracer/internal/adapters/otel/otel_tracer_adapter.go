package otel

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"

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
	ctx            context.Context
	logger         loggercontracts.Logger
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
}

func NewOtelTracerAdapter(
	name string,
	cfg *models.TracerOptions,
	ctx context.Context,
	resource *sdkresource.Resource,
	logger loggercontracts.Logger,
) (contracts.TracerAdapter, error) {
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

func (o *otelTracerAdapter) Start(event string) (context.Context, contracts.TracerSpan) {
	ctx, span := o.tracer.Start(o.ctx, event)
	return ctx, &otelTracerSpan{span: span}
}

func (o *otelTracerAdapter) Close(ctx context.Context) error {
	return o.tracerProvider.Shutdown(ctx)
}

func (o *otelTracerAdapter) Name() string                   { return o.name }
func (o *otelTracerAdapter) Logger() loggercontracts.Logger { return o.logger }
func (o *otelTracerAdapter) Type() tracertype.TracerType {
	return tracertype.TracerTypes.OPEN_TELEMETRY
}
