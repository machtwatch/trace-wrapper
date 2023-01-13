package tracer

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/alvian-machtwatch/tracer/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type SpanTestSuite struct {
	suite.Suite
}

func TestSpanTestSuite(t *testing.T) {
	suite.Run(t, new(SpanTestSuite))
}

func (suite *SpanTestSuite) SetupTest() {
	tc := internal.NewNoopTracerProvider()
	otel.SetTracerProvider(tc)
}

func (a *SpanTestSuite) TestAddEvent() {
	ctx := context.Background()
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	_, span := Span("spanName", SetAttribute("key", "value")).Start(ctx)
	defer span.End()
	span.AddEvent("event")
	a.NotEmpty(out)
}

func (a *SpanTestSuite) TestAddEventWithAttributes() {
	ctx := context.Background()
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	_, span := Span("spanName", SetAttribute("key", "value")).Start(ctx)
	defer span.End()
	span.AddEvent("event", SetAttribute("key", "value"))
	a.NotEmpty(out)
}

func (a *SpanTestSuite) TestLogDebug() {
	ctx := context.Background()
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	_, span := Span("spanName").Start(ctx)
	span.LogDebug("log debug", SetAttribute("key", "value"))
	a.NotEmpty(out)
}

func (a *SpanTestSuite) TestLogError() {
	ctx := context.Background()
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	_, span := Span("spanName").Start(ctx)
	span.LogError("log debug", errors.New("err"), SetAttribute("key", "value"))
	a.NotEmpty(out)
}

func (a *SpanTestSuite) TestSetStatusSuccess() {
	ctx := context.Background()
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	_, span := Span("spanName").Start(ctx)
	span.Success()
	a.NotEmpty(out)
}

func (a *SpanTestSuite) TestSetStatusError() {
	ctx := context.Background()
	out := &bytes.Buffer{}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	_, span := Span("spanName").Start(ctx)
	span.Error(errors.New("string"))
	a.NotEmpty(out)
}
