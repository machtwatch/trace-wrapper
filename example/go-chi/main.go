package main

import (
	"context"
	"log"
	"net/http"

	tracer "github.com/alvian-machtwatch/tracer"
	"github.com/go-chi/chi"
)

func main() {

	config := &tracer.Config{
		PackageName: "github.com/alvian-machtwatch/tracer/example",
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

	//set request id
	r.Use(tracer.RequestID)

	//add tracer to middleware
	r.Use(tp.Middleware)

	//handler
	r.Get("/get", func(w http.ResponseWriter, r *http.Request) {
		tracer.Span("get", tracer.SetAttribute("key", "value")).Start(ctx)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("success"))
	})
	http.ListenAndServe(":3001", r)
}
