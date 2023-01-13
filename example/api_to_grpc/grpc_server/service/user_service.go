package service

import (
	"context"

	tracer "github.com/alvian-machtwatch/tracer"
	"github.com/alvian-machtwatch/tracer/example/api_to_grpc/grpc_server/proto/user"
	"github.com/golang/protobuf/ptypes/empty"
)

type UserServer struct {
	user.UnimplementedUsersServer
}

var localStorage *user.UserList

func init() {
	localStorage = new(user.UserList)
	localStorage.List = make([]*user.User, 0)
}

func (UserServer) Register(ctx context.Context, param *user.User) (*empty.Empty, error) {
	_, span := tracer.Span("user-register").Start(ctx)
	defer span.End()

	localStorage.List = append(localStorage.List, param)

	span.Success()
	span.LogDebug("params", tracer.SetAttribute("params", param.String()))
	return new(empty.Empty), nil
}

func (UserServer) List(ctx context.Context, void *empty.Empty) (*user.UserList, error) {
	_, span := tracer.Span("user-list").Start(ctx)
	defer span.End()
	span.Success()
	span.LogDebug("results", tracer.SetAttribute("results", localStorage.String()))

	return localStorage, nil
}
