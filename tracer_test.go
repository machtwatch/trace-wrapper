package tracer

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type TracerTestSuite struct {
	suite.Suite
}

func TestTracerTestSuite(t *testing.T) {
	suite.Run(t, new(TracerTestSuite))
}

func (a *TracerTestSuite) TestTracerErrorResponse() {
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()
	config := &Config{
		PackageName: "github.com/alvian-machtwatch/logger",
		ServiceName: "Service B",
		Version:     "0.0.1",
		Url:         ts.URL + "/api/traces",
		Environment: "development",
	}
	tc, err := New(config)
	tc.Middleware(handler)
	ctx := context.Background()
	defer func() { tc.Shutdown(ctx) }()
	a.Empty(err)
}
