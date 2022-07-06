package mock

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type MockServerSteam struct {
	Ctx context.Context
}

func (s MockServerSteam) SetHeader(metadata.MD) error { return nil }

func (s MockServerSteam) SendHeader(metadata.MD) error { return nil }

func (s MockServerSteam) SetTrailer(metadata.MD) {}

func (s MockServerSteam) Context() context.Context { return s.Ctx }

func (s MockServerSteam) SendMsg(m interface{}) error { return nil }

func (s MockServerSteam) RecvMsg(m interface{}) error { return nil }
