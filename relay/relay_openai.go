package relay

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model"
)

// TODO: read max_tokens, n, stop, frequency_penalty, presence_penalty from meta

func OpenAIChatCompletionRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
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
	return chatGPT2Grpc(resp)
}

func OpenAIChatCompletionStreamingRequest(rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	fmt.Println("OPENAI_LLM_API")
	config := openai.DefaultConfig(info.GetToken())
	config.BaseURL = info.GetUrl()
	messages := rs.GetMessages()

	client := openai.NewClientWithConfig(config)

	ctx := context.Background()
	var openaiMessages []openai.ChatCompletionMessage
	for _, m := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}

	req := openai.ChatCompletionRequest{
		Model:       info.GetModel(),
		Messages:    openaiMessages,
		TopP:        float32(info.GetTopP()),
		Temperature: float32(info.GetTemperature()),
		Stream:      true,
	}

	resp, err := client.CreateChatCompletionStream(ctx, req)

	if err != nil {
		return err
	}

	go func() {
		defer resp.Close()
		i := 0
		for {
			response, err := resp.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("\nStream finished")
				sv.Send(&model.PartialLLMResponse{
					Response: &model.PartialLLMResponse_Done{},
					LlmTokenCount: &model.LLMTokenCount{
						CompletionToken: int64(i),
					},
				})
				return
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}

			if len(response.Choices) != 0 {
				message := response.Choices[0].Delta.Content
				fmt.Printf("\nStream message %d: %s\n", i, message)
				i++
				pr := model.PartialLLMResponse{
					Response: &model.PartialLLMResponse_Content{
						Content: message,
					},
				}
				sv.Send(&pr)
			}
		}
	}()
	return nil
}
