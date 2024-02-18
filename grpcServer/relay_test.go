package grpcServer_test

import (
	"context"
	"log"
	"os"
	"testing"

	"go.limit.dev/unollm/grpcServer"

	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/utils"

	"github.com/joho/godotenv"
)

func TestOpenAI(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	openaiApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.OPENAI_LLM_API,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.openai-sb.com/v1",
		Token:       openaiApiKey,
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
	t.Log("res: ", res)
}

func TestChatGLM(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.CHATGLM_LLM_API,
		Model:       "glm-3-turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       zhipuaiApiKey,
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

func TestChatGLMStreaming(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.CHATGLM_LLM_API,
		Model:       "glm-3-turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       zhipuaiApiKey,
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

func TestOpenAItreaming(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	openaiApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.OPENAI_LLM_API,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.openai-sb.com/v1",
		Token:       openaiApiKey,
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
		t.Log(res)
		if res.LlmTokenCount != nil {
			t.Log(res.LlmTokenCount)
			return
		}
	}
}
