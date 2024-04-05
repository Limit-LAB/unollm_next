package tests_grpc

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.limit.dev/unollm/grpcServer"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/utils"
)

func TestAnthropic(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	AnthropicApiKey := os.Getenv("TEST_ANTHROPIC_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.ANTHROPIC_LLM_API,
		Model:       "claude-3-opus-20240229",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       AnthropicApiKey,
	}
	req := model.LLMRequestSchema{
		Messages:       messages,
		LlmRequestInfo: &req_info,
	}
	mockServer := grpcServer.UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	log.Println("res: ", res)
}

func TestAnthropicStreaming(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	AnthropicApiKey := os.Getenv("TEST_ANTHROPIC_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.ANTHROPIC_LLM_API,
		Model:       "claude-3-opus-20240229",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       AnthropicApiKey,
	}
	req := model.LLMRequestSchema{
		Messages:       messages,
		LlmRequestInfo: &req_info,
	}
	mockServer := grpcServer.UnoForwardServer{}
	mockServerPipe := utils.MockServerStream{
		Stream: make(chan *model.PartialLLMResponse, 1000),
	}
	err := mockServer.StreamRequestLLM(&req, &mockServerPipe)
	if err != nil {
		t.Error(err)
	}
	for {
		res := <-mockServerPipe.Stream
		log.Println(res)
		if res.LlmTokenCount != nil {
			log.Println(res.LlmTokenCount)
			return
		}
	}
}
