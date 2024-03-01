package utils

import (
	"context"

	"go.limit.dev/unollm/model"
	"google.golang.org/grpc/metadata"
)

type MockServerStream struct {
	Stream  chan *model.PartialLLMResponse
	header  metadata.MD
	trailer metadata.MD
	ctx     context.Context
}

func (m *MockServerStream) Send(res *model.PartialLLMResponse) error {
	m.Stream <- res
	return nil
}

func NewMockServerStream(ctx context.Context) *MockServerStream {
	return &MockServerStream{
		ctx: ctx,
	}
}

func (m *MockServerStream) SetHeader(md metadata.MD) error {
	m.header = md
	return nil
}

func (m *MockServerStream) SendHeader(md metadata.MD) error {
	m.header = md
	return nil
}

func (m *MockServerStream) SetTrailer(md metadata.MD) {
	m.trailer = md
}

func (m *MockServerStream) Context() context.Context {
	return m.ctx
}

func (m *MockServerStream) SendMsg(msg interface{}) error {
	// Mock implementation, no action needed
	return nil
}

func (m *MockServerStream) RecvMsg(msg interface{}) error {
	// Mock implementation, no action needed
	return nil
}
