package auth

import (
	"context"
	"testing"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMetadataKey(t *testing.T) {
	want := metadataKey{}

	got := MetadataKey()

	if got != want {
		t.Errorf("metadata %v; want %v", got, want)
	}
}

func TestDefaultAuth(t *testing.T) {
	wantErr := status.Error(codes.Unauthenticated, "Unauthenticated")

	gotCtx, gotErr := DefaultAuth(context.TODO())

	if gotCtx != nil {
		t.Errorf("context %v; want <nil>", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestNoAuth(t *testing.T) {
	ctx := context.TODO()
	wantCtx := context.WithValue(ctx, metadataKey{}, "")

	gotCtx, gotErr := NoAuth(ctx)

	if gotCtx == nil {
		t.Errorf("context <nil>; want %v", wantCtx)
	}
	gotCtxValue, ok := gotCtx.Value(metadataKey{}).(string)
	if !ok || gotCtxValue != "" {
		t.Errorf("context %v; want %v", gotCtx, wantCtx)
	}
	if gotErr != nil {
		t.Errorf("error %v; want <nil>", gotErr)
	}
}

func TestUserAuth_failureNotBasic(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "bearer worldhello",
	))
	_, wantErr := grpc_auth.AuthFromMD(ctx, "basic")

	gotCtx, gotErr := UserAuth(ctx)

	if gotCtx != nil {
		t.Errorf("context %v; want nil", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestUserAuth_failureWrongToken(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "basic xxx",
	))
	wantErr := status.Error(codes.Unauthenticated, "Unauthenticated")

	gotCtx, gotErr := UserAuth(ctx)

	if gotCtx != nil {
		t.Errorf("context %v; want nil", gotCtx)
	}
	if gotErr == nil || gotErr.Error() != wantErr.Error() {
		t.Errorf("error %v; want %v", gotErr, wantErr)
	}
}

func TestUserAuth_success(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(
		"authorization", "basic worldhello",
	))
	wantCtx := context.WithValue(ctx, metadataKey{}, "")

	gotCtx, gotErr := UserAuth(ctx)

	if gotCtx == nil {
		t.Errorf("context <nil>; want %v", wantCtx)
	}
	gotCtxValue, ok := gotCtx.Value(metadataKey{}).(string)
	if !ok || gotCtxValue != "" {
		t.Errorf("context %v; want %v", gotCtx, wantCtx)
	}
	if gotErr != nil {
		t.Errorf("error %v; want <nil>", gotErr)
	}
}
