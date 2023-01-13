package tracer

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (ts tracerSpan) AddEvent(description string, attributes ...tracerAttribute) {
	var arr []attribute.KeyValue = []attribute.KeyValue{}
	for _, attr := range attributes {
		arr = append(arr, attr.getAttribute())
	}

	ts.span.AddEvent(description, trace.WithAttributes(arr...))
}
