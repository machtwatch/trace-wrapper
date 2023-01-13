package tracer

import (
	"github.com/alvian-machtwatch/tracer/internal"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TracerSpan interface {
	Success()
	Error(err error)
	End()
	AddEvent(description string, attributes ...tracerAttribute)
	LogDebug(description string, attributes ...tracerAttribute)
	LogError(description string, err error, attributes ...tracerAttribute)
}

type tracerSpan struct {
	span       trace.Span
	stackTrace *internal.CallInfo
}

func newTracerSpan(span trace.Span, stackTrace *internal.CallInfo) TracerSpan {
	return tracerSpan{
		span:       span,
		stackTrace: stackTrace,
	}
}

func (ts tracerSpan) Success() {
	ts.span.SetStatus(codes.Ok, "")
}

func (ts tracerSpan) Error(err error) {
	ts.span.SetStatus(codes.Error, err.Error())
	ts.span.RecordError(err)
}

func (ts tracerSpan) End() {
	ts.span.End()
}
