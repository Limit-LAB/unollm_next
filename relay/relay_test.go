package relay_test

import (
	"context"
	"fmt"
	"limit.dev/unollm/model/unoLlmMod"
	"log"
	"os"
	"testing"

	"limit.dev/unollm/relay"

	"github.com/joho/godotenv"
	"google.golang.org/grpc/metadata"
)

func TestOpenAI(t *testing.T) {
	godotenv.Load()

	messages := make([]*unoLlmMod.LLMChatCompletionMessage, 0)
	messages = append(messages, &unoLlmMod.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	openaiApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := unoLlmMod.LLMRequestInfo{
		LlmApiType:  relay.OPENAI_LLM_API,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.openai-sb.com/v1",
		Token:       openaiApiKey,
	}
	req := unoLlmMod.LLMRequestSchema{
		Messages:       messages,
		LlmRequestInfo: &req_info,
	}
	mockServer := relay.UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log("res: ", res)
}

func TestChatGLM(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*unoLlmMod.LLMChatCompletionMessage, 0)
	messages = append(messages, &unoLlmMod.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	req_info := unoLlmMod.LLMRequestInfo{
		LlmApiType:  relay.CHATGLM_LLM_API,
		Model:       "chatglm_turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       zhipuaiApiKey,
	}
	req := unoLlmMod.LLMRequestSchema{
		Messages:       messages,
		LlmRequestInfo: &req_info,
	}
	mockServer := relay.UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log("res: ", res)
}

type MockServerStream struct {
	stream  []*unoLlmMod.PartialLLMResponse
	header  metadata.MD
	trailer metadata.MD
	ctx     context.Context
}

func (m *MockServerStream) Send(res *unoLlmMod.PartialLLMResponse) error {
	fmt.Println(res)
	m.stream = append(m.stream, res)
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

func TestChatGLMStreaming(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	messages := make([]*unoLlmMod.LLMChatCompletionMessage, 0)
	messages = append(messages, &unoLlmMod.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	req_info := unoLlmMod.LLMRequestInfo{
		LlmApiType:  relay.CHATGLM_LLM_API,
		Model:       "chatglm_turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       zhipuaiApiKey,
	}
	req := unoLlmMod.LLMRequestSchema{
		Messages:       messages,
		LlmRequestInfo: &req_info,
	}
	mockServer := relay.UnoForwardServer{}
	mockServerPipe := MockServerStream{
		stream: []*unoLlmMod.PartialLLMResponse{},
	}
	err = mockServer.StreamRequestLLM(&req, &mockServerPipe)
	if err != nil {
		t.Error(err)
	}

}
