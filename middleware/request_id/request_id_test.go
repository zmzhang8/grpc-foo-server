package request_id

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"

	"foo/mock"
)

func TestHeaderKey(t *testing.T) {
	want := headerKey

	got := HeaderKey()

	if got != want {
		t.Errorf("metadata %v; want %v", got, want)
	}
}

func TestMetadataKey(t *testing.T) {
	want := metadataKey{}

	got := MetadataKey()

	if got != want {
		t.Errorf("metadata %v; want %v", got, want)
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	ctx := context.TODO()
	info := grpc.UnaryServerInfo{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "yyy", errors.New("zzz")
	}

	resp, err := UnaryServerInterceptor()(ctx, nil, &info, handler)

	if resp != "yyy" {
		t.Errorf("resp %v; want yyy", resp)
	}
	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	ctx := context.TODO()
	stream := mock.MockServerSteam{Ctx: ctx}
	info := grpc.StreamServerInfo{}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errors.New("zzz")
	}

	err := StreamServerInterceptor()(nil, stream, &info, handler)

	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}
