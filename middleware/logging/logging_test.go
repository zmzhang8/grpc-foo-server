package logging

import (
	"context"
	"errors"
	"os"
	"testing"

	"google.golang.org/grpc"

	"foo/lib/log"
	"foo/mock"
)

func TestMetadataKey(t *testing.T) {
	want := metadataKey{}

	got := MetadataKey()

	if got != want {
		t.Errorf("metadata %v; want %v", got, want)
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	fieldKey := struct{}{}
	fieldName := "key"
	filedValue := "value"
	fields := []Field{
		{MetadataKey: fieldKey, Name: fieldName},
	}
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), fieldKey, filedValue)
	info := grpc.UnaryServerInfo{
		FullMethod: "/test.TestServer/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "yyy", errors.New("zzz")
	}

	resp, err := UnaryServerInterceptor(logger, fields...)(ctx, nil, &info, handler)

	if resp != "yyy" {
		t.Errorf("resp %v; want yyy", resp)
	}
	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	fieldKey := struct{}{}
	fieldName := "key"
	filedValue := "value"
	fields := []Field{
		{MetadataKey: fieldKey, Name: fieldName},
	}
	logger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), fieldKey, filedValue)
	stream := mock.MockServerSteam{Ctx: ctx}
	info := grpc.StreamServerInfo{
		FullMethod: "/test.TestServer/Hello",
	}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errors.New("zzz")
	}

	err := StreamServerInterceptor(logger, fields...)(nil, stream, &info, handler)

	if err.Error() != "zzz" {
		t.Errorf("err %v; want zzz", err)
	}
}
