package tracer

import (
	"context"

	"github.com/alvian-machtwatch/tracer/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TracerSpanManager interface {
	Start(context context.Context) (context.Context, TracerSpan)
}

type tracerSpanManager struct {
	spanName   string
	attributes []tracerAttribute
	span       trace.Span
}

func Span(spanName string, attributes ...tracerAttribute) TracerSpanManager {
	return tracerSpanManager{
		spanName:   spanName,
		attributes: attributes,
	}
}

func (ts tracerSpanManager) Start(context context.Context) (context.Context, TracerSpan) {
	stackTrace := internal.RetrieveCallInfo()

	tracer := otel.Tracer(stackTrace.PackageName)
	context, ts.span = tracer.Start(context, ts.spanName)
	for _, attr := range ts.attributes {
		ts.span.SetAttributes(attr.getAttribute())
	}

	return context, newTracerSpan(ts.span, stackTrace)
}
