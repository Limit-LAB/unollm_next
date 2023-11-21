package relay

import (
	"context"
	"fmt"
	"limit.dev/unollm/model/unoLlmMod"

	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: read max_tokens, n, stop, frequency_penalty, presence_penalty from meta

func OpenaiBlockingRequest(ctx context.Context, rs *unoLlmMod.LLMRequestSchema) (*unoLlmMod.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("OPENAI_LLM_API")
	config := openai.DefaultConfig(info.GetToken())
	config.BaseURL = info.GetUrl()
	client := openai.NewClientWithConfig(config)
	messages := rs.GetMessages()
	var openaiMessages []openai.ChatCompletionMessage
	for _, m := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       info.GetModel(),
			Messages:    openaiMessages,
			TopP:        float32(info.GetTopP()),
			Temperature: float32(info.GetTemperature()),
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return chatGPTTranslateToRelay(resp)
}
