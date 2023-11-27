package respTransformer

import (
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.limit.dev/unollm/model"
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
