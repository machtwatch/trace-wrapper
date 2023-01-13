package tracer

import (
	"context"
	"net/http"

	"github.com/alvian-machtwatch/tracer/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type tracer struct {
	Shutdown func(ctx context.Context) error
	config   *Config
}

func New(config *Config) (*tracer, error) {

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			attribute.String("version", config.Version),
			attribute.String("environment", config.Environment),
			attribute.String("id", config.PackageName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	setGlobalLevel(config)
	tc := &tracer{
		Shutdown: tp.Shutdown,
		config:   config,
	}
	return tc, nil
}

func Mock(config *Config) (*tracer, error) {
	tc := internal.NewNoopTracerProvider()
	otel.SetTracerProvider(tc)
	setGlobalLevel(config)

	shutdown := func(ctx context.Context) error {
		return nil
	}
	tp := &tracer{
		Shutdown: shutdown,
		config:   config,
	}
	return tp, nil
}

func (t tracer) Middleware(next http.Handler) http.Handler {
	return tracerMiddleware{
		Tracer:      otel.GetTracerProvider().Tracer(t.config.PackageName),
		Handler:     next,
		Propagators: otel.GetTextMapPropagator(),
	}
}
