package helloworld

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"foo/lib/auth"
	"foo/lib/custom_context"
	pb "foo/pb_gen/helloworld"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(
	ctx context.Context,
	in *pb.HelloRequest,
) (*pb.HelloReply, error) {
	logger, _ := custom_context.Logger(ctx)
	logger.Infow("Received ", "name", in.GetName())
	if len(in.Name) == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Name cannot be empty",
		)
	}
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *greeterServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return auth.UserAuth(ctx)
}

func NewGreeterServer() *greeterServer {
	return &greeterServer{}
}
