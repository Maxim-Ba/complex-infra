package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func InitTracer(jaegerURL string, serviceName string) (trace.Tracer, error) {
	exporter, err := NewJaegerExporter(jaegerURL)
	if err != nil {
		return nil, fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("initialize provider: %w", err)
	}

	otel.SetTracerProvider(tp) // задает глобальный провайдер !
	Tracer = tp.Tracer("main tracer")
	return Tracer, nil
}
