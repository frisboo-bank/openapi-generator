package otel

import trace "go.opentelemetry.io/otel/trace"

type otelTracerSpan struct {
	span trace.Span
}

func (s *otelTracerSpan) End()                  { s.span.End() }
func (s *otelTracerSpan) RecordError(err error) { s.span.RecordError(err) }
