package tracer

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setGlobalLevel(config *Config) {
	switch config.Environment {
	case "development":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "staging", "production":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		error := errors.New("config.Env is not set or value is not development, staging, or production")
		panic(error)
	}
}

func (ts tracerSpan) LogDebug(description string, attributes ...tracerAttribute) {
	log := ts.prependLog(description, attributes...)
	log.Debug().Msg(description)
}

func (ts tracerSpan) LogError(description string, err error, attributes ...tracerAttribute) {
	log := ts.prependLog(description, attributes...)
	log.Error().Err(err).Msg(description)
}

func (ts tracerSpan) prependLog(description string, attributes ...tracerAttribute) zerolog.Logger {
	sublogger := log.With().Int("line", ts.stackTrace.Line).
		Str("func_name", ts.stackTrace.FuncName).
		Str("filename", ts.stackTrace.FileName).
		Str("package_name", ts.stackTrace.PackageName).
		Str("span_id", ts.span.SpanContext().SpanID().String()).
		Str("trace_id", ts.span.SpanContext().TraceID().String())

	for _, attr := range attributes {
		sublogger.Str(attr.Key, attr.Value)
	}

	return sublogger.Logger()
}
