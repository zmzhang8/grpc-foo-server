package helloworld

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "foo/pb_gen/helloworld"
)

func TestHealthServer_Check_success(t *testing.T) {
	serviceInfo := map[string]grpc.ServiceInfo{
		"test.TestServer.Hello": {},
	}
	s := NewHealthServer(serviceInfo)
	ctx := context.TODO()
	req := pb.HealthCheckRequest{
		Service: "test.TestServer.Hello",
	}
	wantStatus := pb.HealthCheckResponse_SERVING

	resp, err := s.Check(ctx, &req)

	if err != nil {
		t.Errorf("err %v; want <nil>", err)
	}

	if resp.Status != wantStatus {
		t.Errorf("message %v; want %v", resp.Status, wantStatus)
	}
}

func TestHealthServer_Check_failure(t *testing.T) {
	serviceInfo := map[string]grpc.ServiceInfo{}
	s := NewHealthServer(serviceInfo)
	ctx := context.TODO()
	req := pb.HealthCheckRequest{
		Service: "test.TestServer.Hello",
	}
	wantErr := status.Errorf(
		codes.NotFound,
		"Service not found",
	)

	_, err := s.Check(ctx, &req)

	if err == nil {
		t.Errorf("err <nil>; want %v", wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("err %v; want %v", err, wantErr)
	}
}
