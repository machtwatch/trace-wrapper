package tracer

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type MiddlewareTestSuite struct {
	suite.Suite
	config *Config
}

func (suite *MiddlewareTestSuite) SetupTest() {
	config := &Config{
		PackageName: "github.com/alvian-machtwatch/logger",
		ServiceName: "Service B",
		Version:     "0.0.1",
		Url:         "http://localhost:14268/api/traces",
		Environment: "development",
	}
	suite.config = config
	Mock(config)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (a *MiddlewareTestSuite) TestMiddlewareErrorResponse() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	})
	tm := tracerMiddleware{
		Tracer:      otel.Tracer(a.config.PackageName),
		Propagators: otel.GetTextMapPropagator(),
		Handler:     handler,
	}
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	m := chi.NewRouter()
	m.Method("GET", "/foo", handler)
	ts := httptest.NewServer(m)
	defer ts.Close()
	req := httptest.NewRequest("GET", ts.URL+"/foo", nil)
	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	tm.ServeHTTP(httptest.NewRecorder(), req)
	a.NotEmpty(out.String())
}

func (a *MiddlewareTestSuite) TestMiddlewareSuccessResponse() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		w.WriteHeader(http.StatusAccepted)
	})
	tm := tracerMiddleware{
		Tracer:      otel.Tracer(a.config.PackageName),
		Propagators: otel.GetTextMapPropagator(),
		Handler:     handler,
	}
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	m := chi.NewRouter()
	m.Method("GET", "/foo", handler)
	ts := httptest.NewServer(m)
	defer ts.Close()
	req := httptest.NewRequest("GET", ts.URL+"/foo", nil)
	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	tm.ServeHTTP(httptest.NewRecorder(), req)
	a.NotEmpty(out.String())
}
