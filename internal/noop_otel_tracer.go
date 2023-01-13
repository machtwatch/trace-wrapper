package internal

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	trace "go.opentelemetry.io/otel/trace"
)

type traceContextKeyType int

const currentSpanKey traceContextKeyType = iota

// nonRecordingSpan is a minimal implementation of a Span that wraps a
// SpanContext. It performs no operations other than to return the wrapped
// SpanContext.
type nonRecordingSpan struct {
	noopSpan
	SC trace.SpanContext
}

// SpanContext returns the wrapped SpanContext.
func (s nonRecordingSpan) SpanContext() trace.SpanContext {
	config := trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
		Remote:  false,
	}
	s.SC = trace.NewSpanContext(config)
	return s.SC

}

// NewNoopTracerProvider returns an implementation of TracerProvider that
// performs no operations. The Tracer and Spans created from the returned
// TracerProvider also perform no operations.
func NewNoopTracerProvider() trace.TracerProvider {
	return noopTracerProvider{}
}

type noopTracerProvider struct{}

var _ trace.TracerProvider = noopTracerProvider{}

// Tracer returns noop implementation of Tracer.
func (p noopTracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return noopTracer{}
}

// noopTracer is an implementation of Tracer that preforms no operations.
type noopTracer struct {
}

var _ trace.Tracer = noopTracer{}

// Start carries forward a non-recording Span, if one is present in the context, otherwise it
// creates a no-op Span.
func (t noopTracer) Start(ctx context.Context, name string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)
	if _, ok := span.(nonRecordingSpan); !ok {
		// span is likely already a noopSpan, but let's be sure
		span = noopSpan{}
	}
	return trace.ContextWithSpan(ctx, span), span
}

// noopSpan is an implementation of Span that preforms no operations.
type noopSpan struct {
}

var _ trace.Span = noopSpan{}

// SpanContext returns an empty span context.
func (noopSpan) SpanContext() trace.SpanContext {
	config := trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
		Remote:  false,
	}
	return trace.NewSpanContext(config)
}

// IsRecording always returns false.
func (noopSpan) IsRecording() bool { return false }

// SetStatus does nothing.
func (noopSpan) SetStatus(code codes.Code, description string) {
	log.Debug().Str("code", code.String()).Str("description", description).Msg("SetStatus")
}

// SetError does nothing.
func (noopSpan) SetError(value bool) {
	log.Debug().Bool("value", value).Msg("SetError")
}

// SetAttributes does nothing.
func (a noopSpan) SetAttributes(attr ...attribute.KeyValue) {
	log.Debug().Interface("attr", attr).Msg("SetAttributes")

}

// End does nothing.
func (noopSpan) End(...trace.SpanEndOption) {
	log.Debug().Msg("End")
}

// RecordError does nothing.
func (noopSpan) RecordError(error, ...trace.EventOption) {
	log.Debug().Msg("RecordError")
}

// AddEvent does nothing.
func (noopSpan) AddEvent(string, ...trace.EventOption) {
	log.Debug().Msg("AddEvent")
}

// SetName does nothing.
func (noopSpan) SetName(name string) {
	log.Debug().Str("name", name).Msg("SetName")

}

// TracerProvider returns a no-op TracerProvider.
func (noopSpan) TracerProvider() trace.TracerProvider { return noopTracerProvider{} }
