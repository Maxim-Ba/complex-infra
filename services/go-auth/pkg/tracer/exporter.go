package tracer

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)
type SpanExporter interface {
    ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error
    Shutdown(ctx context.Context) error
}
// NewJaegerExporter creates new jaeger exporter
//
func NewJaegerExporter(url string) (tracesdk.SpanExporter, error) {
	return otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(url), // "jaeger:4318"
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithURLPath("/v1/traces"), // Это важно для HTTP
	)
}
