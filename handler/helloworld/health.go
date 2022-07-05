package helloworld

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"foo/lib/auth"
	pb "foo/pb_gen/helloworld"
)

type healthServer struct {
	pb.UnimplementedHealthServer
	services map[string]grpc.ServiceInfo
}

func (s *healthServer) Check(
	ctx context.Context,
	in *pb.HealthCheckRequest,
) (*pb.HealthCheckResponse, error) {
	if in.Service != "" {
		if _, ok := s.services[in.Service]; !ok {
			return nil, status.Errorf(
				codes.NotFound,
				"Service not found",
			)
		}
	}

	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(
	in *pb.HealthCheckRequest,
	stream pb.Health_WatchServer,
) error {
	if in.Service != "" {
		if _, ok := s.services[in.Service]; !ok {
			return status.Errorf(
				codes.NotFound,
				"Service not found",
			)
		}
	}

	stream.Send(&pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING})
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			stream.Send(&pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING})
		}
	}
}

func (s *healthServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return auth.NoAuth(ctx)
}

func NewHealthServer(services map[string]grpc.ServiceInfo) *healthServer {
	return &healthServer{
		services: services,
	}
}
