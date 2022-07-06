package helloworld

import (
	"context"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"foo/lib/log"
	"foo/middleware/logging"
	pb "foo/pb_gen/helloworld"
)

func TestGreeterServer_SayHello_success(t *testing.T) {
	s := NewGreeterServer()
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), logging.MetadataKey(), logger)
	req := pb.HelloRequest{
		Name: "world",
	}
	wantMessage := "Hello world"

	resp, err := s.SayHello(ctx, &req)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}

	if resp.Message != wantMessage {
		t.Errorf("message %v; want %v", resp.Message, wantMessage)
	}
}

func TestGreeterServer_SayHello_failure(t *testing.T) {
	s := NewGreeterServer()
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), logging.MetadataKey(), logger)
	req := pb.HelloRequest{
		Name: "",
	}
	wantErr := status.Errorf(
		codes.InvalidArgument,
		"Name cannot be empty",
	)

	_, err := s.SayHello(ctx, &req)

	if err == nil {
		t.Errorf("err <nil>; want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("err %v; want %v", err, wantErr)
	}
}
