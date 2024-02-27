package grpcServer_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.limit.dev/unollm/grpcServer"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/utils"
)

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
		Model:       ChatGLM.ModelGLM3Turbo,
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
		Model:       ChatGLM.ModelGLM3Turbo,
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

func TestChatGLMFunctionCalling(t *testing.T) {
	godotenv.Load("../.env")

	messages := make([]*model.LLMChatCompletionMessage, 0)
	messages = append(messages, &model.LLMChatCompletionMessage{
		Role:    "user",
		Content: "北京现在什么天气？",
	})
	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	req_info := model.LLMRequestInfo{
		LlmApiType:  grpcServer.CHATGLM_LLM_API,
		Model:       ChatGLM.ModelGLM3Turbo,
		Temperature: 0.9,
		TopP:        0.9,
		TopK:        1,
		Url:         "",
		Token:       zhipuaiApiKey,
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
		t.Error(err)
	}
	log.Printf("res: %#v\n", res.ToolCalls[0])
}
