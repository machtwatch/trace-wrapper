package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	tracer "github.com/alvian-machtwatch/tracer"
	"github.com/alvian-machtwatch/tracer/example/api_to_grpc/grpc_server/proto/user"
	grpcTracer "github.com/alvian-machtwatch/tracer/grpc"
	"github.com/go-chi/chi"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	config := &tracer.Config{
		PackageName: "github.com/alvian-machtwatch/logger",
		ServiceName: "API A",
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
	r.Use(tracer.RequestID)
	r.Use(tp.Middleware)
	r.Get("/get", func(w http.ResponseWriter, r *http.Request) {

		ctx, span := tracer.Span("child1").Start(r.Context())
		defer span.End()
		span.Error(errors.New("err"))
		service := userService()
		user1 := &user.User{
			Id:       "n001",
			Name:     "Noval Agung",
			Password: "kw8d hl12/3m,a",
			Gender:   user.UserGender_MALE,
		}
		service.Register(ctx, user1)
		service.Register(ctx, user1)

		res1, err := service.List(ctx, new(empty.Empty))
		if err != nil {
			log.Fatal(err.Error())
		}
		res1String, _ := json.Marshal(res1.List)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(res1String))
	})
	http.ListenAndServe(":3001", r)
}

func userService() user.UsersClient {
	port := ":3002"
	conn, err := grpcTracer.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return user.NewUsersClient(conn)
}
