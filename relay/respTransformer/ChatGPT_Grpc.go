package respTransformer

import (
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

func ChatGPTToGrpcCompletion(resp openai.ChatCompletionResponse) (*model.LLMResponseSchema, error) {
	if len(resp.Choices) == 0 {
		return nil, status.Errorf(codes.Internal, "OpenAI choices is empty")
	}
	message := resp.Choices[0].Message
	retMessage := model.LLMChatCompletionMessage{
		Role:    message.Role,
		Content: message.Content,
	}
	count := model.LLMTokenCount{
		TotalToken:      int64(resp.Usage.TotalTokens),
		PromptToken:     int64(resp.Usage.PromptTokens),
		CompletionToken: int64(resp.Usage.CompletionTokens),
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}

func ChatGPTToGrpcStream(resp *openai.ChatCompletionStream, sv model.UnoLLMv1_StreamRequestLLMServer) error {
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
			return nil
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return err
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
}
