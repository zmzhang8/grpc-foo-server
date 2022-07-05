package custom_context

import (
	"context"

	"foo/lib/log"
	"foo/middleware/logging"
	"foo/middleware/request_id"
)

func Logger(ctx context.Context) (log.Logger, bool) {
	logger, ok := ctx.Value(logging.MetadataKey()).(log.Logger)
	return logger, ok
}

func RequestId(ctx context.Context) (string, bool) {
	requestId, ok := ctx.Value(request_id.MetadataKey()).(string)
	return requestId, ok
}
