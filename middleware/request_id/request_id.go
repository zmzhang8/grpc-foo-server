package request_id

import (
	"context"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const headerKey = "request-id"

func HeaderKey() string {
	return headerKey
}

type metadataKey struct{}

func MetadataKey() metadataKey {
	return metadataKey{}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestId := uuid.NewString()
		newCtx := context.WithValue(ctx, metadataKey{}, requestId)
		grpc.SetHeader(newCtx, metadata.Pairs(
			headerKey, requestId,
		))

		return handler(newCtx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		requestId := uuid.NewString()
		newCtx := context.WithValue(stream.Context(), metadataKey{}, requestId)
		grpc.SetHeader(newCtx, metadata.Pairs(
			headerKey, requestId,
		))
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
