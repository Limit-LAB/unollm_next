package tests_http_test

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	tests_http "go.limit.dev/unollm/tests_http"
)

func TestMoonShotStreaming(t *testing.T) {
	godotenv.Load("../.env")

	client := tests_http.GetClient(os.Getenv("TEST_MOONSHOT_API"))

	resp, err := client.CreateChatCompletionStream(context.Background(),
		openai.ChatCompletionRequest{
			Model: "moonshot-v1-8k",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: "如果今天下雨，我需要打伞吗？",
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	for {
		cv, e := resp.Recv()
		if e != nil {
			if errors.Is(e, io.EOF) {
				break
			}
			t.Error(e)
			break
		}
		log.Printf("%#v", cv.Choices[0].Delta)
	}
}

func TestMoonShotBlocking(t *testing.T) {
	godotenv.Load("../.env")

	client := tests_http.GetClient(os.Getenv("TEST_MOONSHOT_API"))

	resp, err := client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model: "moonshot-v1-8k",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: "如果今天下雨，我需要打伞吗？",
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#v", resp.Choices[0])
}

func TestMoonShotFunctionCalling(t *testing.T) {
	godotenv.Load("../.env")

	client := tests_http.GetClient(os.Getenv("TEST_MOONSHOT_API"))

	resp, err := client.CreateChatCompletionStream(context.Background(),
		openai.ChatCompletionRequest{
			Model: "moonshot-v1-8k",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: "请返回北京和天津的天气情况。",
				},
			},
			ToolChoice: "auto",
			Tools: []openai.Tool{
				{
					Type: openai.ToolType("function"),
					Function: &openai.FunctionDefinition{
						Name:        "get_weather",
						Description: "Get the weather of a location",
						Parameters: map[string]any{
							"type": "object",
							"properties": map[string]any{
								"location": map[string]any{
									"type":        "string",
									"description": "The city and state, e.g. San Francisco, CA",
								},
								"unit": map[string]any{
									"type": "string",
									"enum": []string{"celsius", "fahrenheit"},
								},
							},
							"required": []string{"location", "unit"},
						},
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	for {
		cv, e := resp.Recv()
		if e != nil {
			if errors.Is(e, io.EOF) {
				break
			}
			log.Panic(e)
			break
		}
		log.Printf("%#v", cv.Choices[0].Delta.ToolCalls)
	}
}
