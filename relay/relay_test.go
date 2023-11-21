package relay_test

import (
	"context"
	"limit.dev/unollm/model"
	"limit.dev/unollm/utils"
	"os"
	"testing"

	"limit.dev/unollm/relay"

	"github.com/joho/godotenv"
)

func TestOpenAI(t *testing.T) {
	godotenv.Load()

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	openaiApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  relay.OPENAI_LLM_API,
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
	mockServer := relay.UnoForwardServer{}
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
		LlmApiType:  relay.CHATGLM_LLM_API,
		Model:       "chatglm_turbo",
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
	mockServer := relay.UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log("res: ", res)
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
		LlmApiType:  relay.CHATGLM_LLM_API,
		Model:       "chatglm_turbo",
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
	mockServer := relay.UnoForwardServer{}
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
