package grpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
)

func NewServer(opt ...grpc.ServerOption) *grpc.Server {
	if opt == nil {
		opt = []grpc.ServerOption{}
	}
	opt = append(opt, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	opt = append(opt, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	return grpc.NewServer(opt...)
}

func Dial(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	if opt == nil {
		opt = []grpc.DialOption{}
	}
	opt = append(opt, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
	opt = append(opt, grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return grpc.Dial(target, opt...)
}
