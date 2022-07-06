package logging

import (
	"context"
	"path"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"foo/lib/log"
)

type metadataKey struct{}

func MetadataKey() metadataKey {
	return metadataKey{}
}

type Field struct {
	MetadataKey interface{}
	Name        string
}

func UnaryServerInterceptor(logger log.Logger, fields ...Field) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now().UTC()
		contextLogger := loggerWithFields(logger, ctx, fields)
		newCtx := context.WithValue(ctx, metadataKey{}, contextLogger)

		contextLogger.Info("Started unary call")

		resp, err := handler(newCtx, req)

		code := status.Code(err)
		summary := summaryArgs(code, info.FullMethod, startTime)
		logwCodeToLevel(contextLogger, code, "Finished unary call", summary...)

		return resp, err
	}
}

func StreamServerInterceptor(logger log.Logger, fields ...Field) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now().UTC()
		ctx := stream.Context()
		contextLogger := loggerWithFields(logger, ctx, fields)
		newCtx := context.WithValue(ctx, metadataKey{}, contextLogger)
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		contextLogger.Info("Started stream call")

		err := handler(srv, wrapped)

		code := status.Code(err)
		summaryArgs := summaryArgs(code, info.FullMethod, startTime)
		logwCodeToLevel(contextLogger, code, "Finished stream call", summaryArgs...)

		return err
	}
}

func loggerWithFields(logger log.Logger, ctx context.Context, fields []Field) log.Logger {
	args := make([]interface{}, 0)
	for _, field := range fields {
		value := ctx.Value(field.MetadataKey)
		args = append(args, field.Name)
		args = append(args, value)
	}

	return logger.With(args...)
}

func logwCodeToLevel(
	logger log.Logger,
	code codes.Code,
	msg string,
	keysAndValues ...interface{},
) {
	switch code {
	case codes.OK:
		logger.Infow(msg, keysAndValues...)
	case codes.Canceled:
		logger.Infow(msg, keysAndValues...)
	case codes.Unknown:
		logger.Errorw(msg, keysAndValues...)
	case codes.InvalidArgument:
		logger.Infow(msg, keysAndValues...)
	case codes.DeadlineExceeded:
		logger.Warnw(msg, keysAndValues...)
	case codes.NotFound:
		logger.Infow(msg, keysAndValues...)
	case codes.AlreadyExists:
		logger.Infow(msg, keysAndValues...)
	case codes.PermissionDenied:
		logger.Warnw(msg, keysAndValues...)
	case codes.Unauthenticated:
		logger.Infow(msg, keysAndValues...) // unauthenticated requests can happen
	case codes.ResourceExhausted:
		logger.Warnw(msg, keysAndValues...)
	case codes.FailedPrecondition:
		logger.Warnw(msg, keysAndValues...)
	case codes.Aborted:
		logger.Warnw(msg, keysAndValues...)
	case codes.OutOfRange:
		logger.Warnw(msg, keysAndValues...)
	case codes.Unimplemented:
		logger.Errorw(msg, keysAndValues...)
	case codes.Internal:
		logger.Errorw(msg, keysAndValues...)
	case codes.Unavailable:
		logger.Warnw(msg, keysAndValues...)
	case codes.DataLoss:
		logger.Errorw(msg, keysAndValues...)
	default:
		logger.Errorw(msg, keysAndValues...)
	}
}

func summaryArgs(
	code codes.Code,
	fullMethod string,
	startTime time.Time,
) []interface{} {
	service := path.Dir(fullMethod)[1:]
	method := path.Base(fullMethod)
	duration := time.Since(startTime)

	return []interface{}{
		"grpc.service", service,
		"grpc.method", method,
		"grpc.code", code,
		"grpc.start_time", startTime,
		"grpc.duration_ms", float64(duration) / float64(time.Millisecond),
	}
}
