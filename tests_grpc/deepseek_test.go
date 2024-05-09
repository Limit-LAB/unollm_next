package tests_grpc_test

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

func TestDeepSeek(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	OPENAIApiKey := os.Getenv("TEST_DEEPSEEK_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.DEEPSEEK_LLM_API,
		Model:       "deepseek-chat",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.deepseek.com/v1",
		Token:       OPENAIApiKey,
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

func TestDeepSeekStreaming(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	OPENAIApiKey := os.Getenv("TEST_DEEPSEEK_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.DEEPSEEK_LLM_API,
		Model:       "deepseek-chat",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.deepseek.com/v1",
		Token:       OPENAIApiKey,
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
		t.Fatal(err)
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

func TestDeepSeekFunctionCalling(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "whats the weather like in Poston?",
	})
	OPENAIApiKey := os.Getenv("TEST_DEEPSEEK_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.DEEPSEEK_LLM_API,
		Model:       "deepseek-chat",
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.deepseek.com/v1",
		Token:       OPENAIApiKey,
		Functions: []*model.Function{
			{
				Name:        "get_weather",
				Description: "Get the weather of a location",
				Parameters: []*model.FunctionCallingParameter{
					{
						Name:        "location",
						Type:        "string",
						Description: "The city and state, e.g. San Francisco, CA",
					},
					{
						Name:  "unit",
						Type:  "string",
						Enums: []string{"celsius", "fahrenheit"},
					},
				},
				Requireds: []string{"location", "unit"},
			},
		},
		UseFunctionCalling: true,
	}
	req := model.LLMRequestSchema{
		Messages:       messages,
		LlmRequestInfo: &req_info,
	}
	mockServer := grpcServer.UnoForwardServer{}
	res, err := mockServer.BlockingRequestLLM(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("res: %#v", res.ToolCalls[0])
}
