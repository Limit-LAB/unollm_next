package tests_grpc_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/grpcServer"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/utils"
)

func TestOpenAI(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	OPENAIApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.OPENAI_LLM_API,
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.openai-sb.com/v1",
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

func TestOpenAIStreaming(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "假如今天下大雨，我是否需要带伞？",
	})
	OPENAIApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.OPENAI_LLM_API,
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.openai-sb.com/v1",
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

func TestOpenAIFunctionCalling(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "whats the weather like in Poston?",
	})
	OPENAIApiKey := os.Getenv("TEST_OPENAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.OPENAI_LLM_API,
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "https://api.openai-sb.com/v1",
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
				Requireds: []string{"location"},
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

func TestOpenAIEmbedding(t *testing.T) {
	godotenv.Load("../.env")
	OPENAIApiKey := os.Getenv("TEST_OPENAI_API")
	mockServer := grpcServer.UnoEmbeddingForwardServer{}
	res, err := mockServer.EmbeddingRequestLLM(context.Background(), &model.EmbeddingRequest{
		EmbeddingRequestInfo: &model.EmbeddingRequestInfo{
			LlmApiType: grpcServer.OPENAI_LLM_API,
			Model:      string(openai.AdaEmbeddingV2),
			Url:        "https://api.openai-sb.com/v1",
			Token:      OPENAIApiKey,
		},
		Text: "我阐述你的梦",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("res: %#v", res)
}
