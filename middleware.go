package tracer

import (
	"net/http"
	"sync"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerMiddleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
type tracerMiddleware struct {
	Tracer      trace.Tracer
	Propagators propagation.TextMapPropagator
	Handler     http.Handler
}

type recordingResponseWriter struct {
	writer  http.ResponseWriter
	written bool
	status  int
}

func (tm tracerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	spanName := r.RequestURI
	opts := []trace.SpanStartOption{
		trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
		trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
		trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(r.RemoteAddr, r.RequestURI, r)...),
		trace.WithSpanKind(trace.SpanKindServer),
	}
	ctx, span := tm.Tracer.Start(ctx, spanName, opts...)
	defer span.End()
	r2 := r.WithContext(ctx)
	rrw := getRRW(w)
	defer putRRW(rrw)
	tm.Handler.ServeHTTP(rrw.writer, r2)

	attrs := semconv.HTTPAttributesFromHTTPStatusCode(rrw.status)
	spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(rrw.status, trace.SpanKindServer)
	span.SetAttributes(attrs...)
	span.SetStatus(spanStatus, spanMessage)
	sublogger := log.With().Timestamp().Str("trace_id", span.SpanContext().TraceID().String()).Str("span_id", span.SpanContext().SpanID().String()).Fields(map[string]interface{}{
		"remote_ip":  r.RemoteAddr,
		"host":       r.Host,
		"proto":      r.Proto,
		"method":     r.Method,
		"path":       r.RequestURI,
		"user_agent": r.Header.Get("User-Agent"),
		"status":     rrw.status,
		"bytes_in":   r.Header.Get("Content-Length"),
		"response":   spanMessage,
	}).Logger()
	code, _ := semconv.SpanStatusFromHTTPStatusCode(rrw.status)
	if code == codes.Error {
		sublogger.Error().Msg("handled_request")
	} else {
		sublogger.Debug().Msg("handled_request")
	}

}

var rrwPool = &sync.Pool{
	New: func() interface{} {
		return &recordingResponseWriter{}
	},
}

func getRRW(writer http.ResponseWriter) *recordingResponseWriter {
	rrw := rrwPool.Get().(*recordingResponseWriter)
	rrw.written = false
	rrw.status = http.StatusOK
	rrw.writer = httpsnoop.Wrap(writer, httpsnoop.Hooks{
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				if !rrw.written {
					rrw.written = true
				}
				return next(b)
			}
		},
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(statusCode int) {
				if !rrw.written {
					rrw.written = true
					rrw.status = statusCode
				}
				next(statusCode)
			}
		},
	})
	return rrw
}

func putRRW(rrw *recordingResponseWriter) {
	rrw.writer = nil
	rrwPool.Put(rrw)
}
