package auth

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type metadataKey struct{}

func MetadataKey() metadataKey {
	return metadataKey{}
}

func DefaultAuth(ctx context.Context) (context.Context, error) {
	return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
}

func NoAuth(ctx context.Context) (context.Context, error) {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/auth#pkg-types
	// The `Context` returned must be a child `Context` of the one passed in
	newCtx := context.WithValue(ctx, metadataKey{}, "")
	return newCtx, nil
}

// Expected header
// key: authorization
// value: basic worldhello
func UserAuth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "basic")
	if err != nil {
		return nil, err
	}

	if token != "worldhello" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/auth#pkg-types
	// The `Context` returned must be a child `Context` of the one passed in
	newCtx := context.WithValue(ctx, metadataKey{}, "")
	return newCtx, nil
}
