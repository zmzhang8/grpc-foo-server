package custom_context

import (
	"context"
	"os"
	"testing"

	"foo/lib/log"
	"foo/middleware/logging"
	"foo/middleware/request_id"
)

func TestLogger_success(t *testing.T) {
	wantLogger := log.NewLogger(log.NewCore(false, os.Stdout, false))
	ctx := context.WithValue(context.TODO(), logging.MetadataKey(), wantLogger)

	gotLogger, gotOk := Logger(ctx)

	if !gotOk {
		t.Errorf("ok %v; want %v", gotOk, true)
	}
	if gotLogger == nil {
		t.Errorf("logger <nil>; want %v", wantLogger)
	}
}

func TestLogger_failure(t *testing.T) {
	_, gotOk := Logger(context.TODO())

	if gotOk {
		t.Errorf("ok %v; want %v", gotOk, false)
	}
}

func TestRequestId_success(t *testing.T) {
	requestId := "xxx"
	ctx := context.WithValue(context.TODO(), request_id.MetadataKey(), requestId)

	gotRequestId, gotOk := RequestId(ctx)

	if !gotOk {
		t.Errorf("ok %v; want %v", gotOk, true)
	}
	if gotRequestId != requestId {
		t.Errorf("request id %v; want %v", gotRequestId, requestId)
	}
}

func TestRequestId_failure(t *testing.T) {
	_, gotOk := RequestId(context.TODO())

	if gotOk {
		t.Errorf("ok %v; want %v", gotOk, false)
	}
}
