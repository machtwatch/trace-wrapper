package main

import (
	"context"
	"log"
	"net"

	tracer "github.com/alvian-machtwatch/tracer"
	"github.com/alvian-machtwatch/tracer/example/api_to_grpc/grpc_server/proto/user"
	"github.com/alvian-machtwatch/tracer/example/api_to_grpc/grpc_server/service"
	grpcTracer "github.com/alvian-machtwatch/tracer/grpc"
)

func main() {
	config := &tracer.Config{
		PackageName: "github.com/alvian-machtwatch/logger",
		ServiceName: "GRPC A",
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
	srv := grpcTracer.NewServer()
	var userSrv service.UserServer
	user.RegisterUsersServer(srv, userSrv)

	l, err := net.Listen("tcp", ":3002")
	if err != nil {
		log.Fatalf("could not listen to %s: %v", ":3002", err)
	}
	if err := srv.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
