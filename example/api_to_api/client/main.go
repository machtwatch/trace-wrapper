package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	tracer "github.com/alvian-machtwatch/tracer"
	"github.com/alvian-machtwatch/tracer/api"
	"github.com/go-chi/chi"
)

func main() {

	config := &tracer.Config{
		PackageName: "github.com/alvian-machtwatch/logger",
		ServiceName: "Service B",
		Version:     "0.0.1",
		Url:         "http://localhost:14268/api/traces",
		Environment: "development",
	}
	tp, err := tracer.New(config)
	ctx := context.Background()

	defer func() { _ = tp.Shutdown(ctx) }()
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	// r.Use(tracer.RequestID)
	// r.Use(tp.Middleware)
	r.Get("/get", api.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := methodService(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("error"))
		}
	}, "get"))

	http.ListenAndServe(":3003", r)
}

func methodService(context context.Context) error {
	ctx, span := tracer.Span("test").Start(context)
	defer span.End()
	request, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/get", nil)

	res, err := tracer.HTTPClientTransporter(http.DefaultTransport).RoundTrip(request)
	if err != nil {
		span.Error(err)
		span.LogError("dummy", err, tracer.SetAttribute("key", "value"))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		err := fmt.Errorf("unexpected response code")
		span.LogError("dummy", err, tracer.SetAttribute("key", "value"))
		return err
	}
	return nil
}
