package respTransformer

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	toolCalls := make([]*model.ToolCall, len(resp.Choices[0].Message.ToolCalls))
	for i, _ := range resp.Choices[0].Message.ToolCalls {
		toolcall := model.ToolCall{
			Id:        resp.Choices[0].Message.ToolCalls[i].ID,
			Name:      resp.Choices[0].Message.ToolCalls[i].Function.Name,
			Arguments: resp.Choices[0].Message.ToolCalls[i].Function.Arguments,
		}
		toolCalls[i] = &toolcall
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
		ToolCalls:     toolCalls,
	}
	return &retResp, nil
}

func ChatGPTToGrpcStream(promptTokens int, resp *openai.ChatCompletionStream, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	defer resp.Close()
	i := 0
	for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("\nStream finished")
			sv.Send(&model.PartialLLMResponse{
				Response: &model.PartialLLMResponse_Done{},
				LlmTokenCount: &model.LLMTokenCount{
					PromptToken:     int64(promptTokens),
					CompletionToken: int64(i),
					TotalToken:      int64(promptTokens + i),
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

			toolCalls := make([]*model.ToolCall, len(response.Choices[0].Delta.ToolCalls))
			for i, _ := range response.Choices[0].Delta.ToolCalls {
				toolcall := model.ToolCall{
					Id:        response.Choices[0].Delta.ToolCalls[i].ID,
					Name:      response.Choices[0].Delta.ToolCalls[i].Function.Name,
					Arguments: response.Choices[0].Delta.ToolCalls[i].Function.Arguments,
				}
				toolCalls[i] = &toolcall
			}
			pr := model.PartialLLMResponse{
				Response: &model.PartialLLMResponse_Content{
					Content: message,
				},
			}
			sv.Send(&pr)
		}
	}
}
