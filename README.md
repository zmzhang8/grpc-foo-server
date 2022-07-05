# Foo Service

## Prerequisites

1. [Go 1.17+](https://go.dev/dl)
1. [Protocol buffer compiler](https://grpc.io/docs/protoc-installation)
1. Go plugins for the protocol buffer compiler
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Quick Start

Start server
```
go run main.go
```

## Generate Client and Server Code

This application leverages Protocol Buffers for service definitions and data serialization. In the "grpc-foo-proto" repository, run the following command to generate gRPC client and server interfaces.
```
protoc --go_opt=paths=source_relative --go_out=/path/to/grpc-foo-server/pb_gen \
  --go-grpc_opt=paths=source_relative --go-grpc_out=/path/to/grpc-foo-server/pb_gen \
  $(find . -iname "*.proto")
```

## Build

1. Build server binary
```
env GOOS=linux GOARCH=amd64 go build -o build/server main.go
```
2. Build Docker image
```
docker build -t grpc-foo-server .
```

## Reference

- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers/docs/overview)
- [Protocol Buffers in Go](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
- [gRPC Documentation](https://grpc.io/docs)
- [gRPC in Go](https://grpc.io/docs/languages/go)
- [gRPC Error Handling](https://www.grpc.io/docs/guides/error)
- [gRPC Performance Best Practices](https://www.grpc.io/docs/guides/performance)
- [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
- [gRPC Metadata in Go](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md)