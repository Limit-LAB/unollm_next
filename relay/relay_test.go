package relay

import (
	"context"
	"limit.dev/unollm/model"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestOpenAI(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	openaiApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  OPENAI_LLM_API,
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
	mockServer := UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log("res: ", res)
}

func TestChatGLM(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  CHATGLM_LLM_API,
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
	mockServer := UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	t.Log("res: ", res)
}
