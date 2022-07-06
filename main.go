package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"google.golang.org/grpc"

	helloworld_handler "foo/handler/helloworld"
	"foo/lib/auth"
	"foo/lib/log"
	mid_logging "foo/middleware/logging"
	mid_request_id "foo/middleware/request_id"
	helloword_pb "foo/pb_gen/helloworld"
)

func main() {
	var (
		debug = flag.Bool("debug", false, "Debug mode")
		port  = flag.Int("port", 8080, "Listen port")
	)
	flag.Parse()

	log.InitDefaultLogger(log.NewCore(false, os.Stdout, *debug))
	defer log.Sync()
	if *debug {
		log.Info("Debug Mode")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	loggingFields := mid_logging.Field{
		MetadataKey: mid_request_id.MetadataKey(),
		Name:        mid_request_id.HeaderKey(),
	}
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			mid_request_id.StreamServerInterceptor(),
			mid_logging.StreamServerInterceptor(log.DefaultLogger, loggingFields),
			grpc_auth.StreamServerInterceptor(auth.DefaultAuth),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			mid_request_id.UnaryServerInterceptor(),
			mid_logging.UnaryServerInterceptor(log.DefaultLogger, loggingFields),
			grpc_auth.UnaryServerInterceptor(auth.DefaultAuth),
		)),
	)
	helloword_pb.RegisterGreeterServer(server, helloworld_handler.NewGreeterServer())
	helloword_pb.RegisterRouteGuideServer(server, helloworld_handler.NewRouteGuideServer())
	helloword_pb.RegisterAccountServer(server, helloworld_handler.NewAccountServer())
	// HealthServer must be registered last
	helloword_pb.RegisterHealthServer(server, helloworld_handler.NewHealthServer(server.GetServiceInfo()))

	log.Info("server listening at", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
