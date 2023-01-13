package tracer

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func HTTPClientTransporter(rt http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(rt)
}
