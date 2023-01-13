package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	tracer "github.com/alvian-machtwatch/tracer"
	"github.com/alvian-machtwatch/tracer/api"
	"github.com/go-chi/chi"
)

func main() {

	config := &tracer.Config{
		PackageName: "github.com/alvian-machtwatch/tracer/example",
		ServiceName: "Service A",
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
		_, span := tracer.Span("child1").Start(r.Context())
		defer span.End()
		span.Error(errors.New("err"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("masuk"))
	}, "get"))
	http.ListenAndServe(":3000", r)
}
