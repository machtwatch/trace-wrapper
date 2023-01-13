package main

import (
	"context"

	tracer "github.com/alvian-machtwatch/tracer"
)

func main() {
	config := &tracer.Config{
		PackageName: "github.com/alvian-machtwatch/tracer/example",
		ServiceName: "Service B",
		Version:     "0.0.1",
		Url:         "http://localhost:14268/api/traces",
		Environment: "development",
	}
	ctx := context.Background()
	tc, _ := tracer.New(config)
	defer tc.Shutdown(ctx)

	_, span := tracer.Span("span").Start(ctx)
	defer span.End()
}
