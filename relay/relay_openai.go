package relay

import (
	"context"
	"fmt"
	"limit.dev/unollm/model"

	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: read max_tokens, n, stop, frequency_penalty, presence_penalty from meta

func OpenaiBlockingRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("OPENAI_LLM_API")
	config := openai.DefaultConfig(info.GetToken())
	config.BaseURL = info.GetUrl()
	client := openai.NewClientWithConfig(config)
	messages := rs.GetMessages()
	openaiMessages := make([]openai.ChatCompletionMessage, 0)
	for _, m := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    info.GetModel(),
			Messages: openaiMessages,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return chatGPTTranslateToRelay(resp)
}
