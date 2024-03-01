package respTransformer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/utils"
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

func GrpcToChatGPTCompletion(model string, resp *model.LLMResponseSchema) openai.ChatCompletionResponse {
	toolCalls := make([]openai.ToolCall, len(resp.ToolCalls))
	for i, toolCall := range resp.ToolCalls {
		toolCalls[i] = openai.ToolCall{
			ID:   toolCall.Id,
			Type: openai.ToolType("function"),
			Function: openai.FunctionCall{
				Name:      toolCall.Name,
				Arguments: toolCall.Arguments,
			},
		}
	}

	message := openai.ChatCompletionMessage{
		Role:      resp.Message.Role,
		Content:   resp.Message.Content,
		ToolCalls: toolCalls,
	}

	return openai.ChatCompletionResponse{
		ID:      "grpc-response",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []openai.ChatCompletionChoice{
			{
				Index:   0,
				Message: message,
			},
		},
	}
}

func OpenAIToGrpcEmbedding(res openai.EmbeddingResponse) (*model.EmbeddingResponse, error) {
	return &model.EmbeddingResponse{
		Dimension: int32(len(res.Data[0].Embedding)),
		Vectors:   res.Data[0].Embedding,
		Usage: &model.LLMTokenCount{
			TotalToken:  int64(res.Usage.TotalTokens),
			PromptToken: int64(res.Usage.PromptTokens),
		},
	}, nil
}

func GrpcStreamToChatGPT(c *gin.Context, model string, sv chan *model.PartialLLMResponse) {
	utils.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		pr := <-sv
		if pr.LlmTokenCount == nil {
			toolcalls := make([]openai.ToolCall, len(pr.GetToolCalls()))
			for i, _ := range pr.GetToolCalls() {
				toolcalls[i] = openai.ToolCall{
					ID:   pr.GetToolCalls()[i].Id,
					Type: openai.ToolType("function"),
					Function: openai.FunctionCall{
						Name:      pr.GetToolCalls()[i].Name,
						Arguments: pr.GetToolCalls()[i].Arguments,
					},
				}
			}
			message := openai.ChatCompletionStreamResponse{
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   model,
				Choices: []openai.ChatCompletionStreamChoice{
					{
						Index: 0,
						Delta: openai.ChatCompletionStreamChoiceDelta{
							Role:      "assistant",
							Content:   pr.GetContent(),
							ToolCalls: toolcalls,
						},
					},
				},
			}
			jsonResponse, _ := json.Marshal(message)
			c.Render(-1, utils.CustomEvent{Data: "data: " + string(jsonResponse)})
			return true
		}

		c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
		return false
	})
}
