package relay

import (
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model"
)

var _ ResponseGrpcTransformer = ChatGPTTranslateToRelay

func ChatGPTTranslateToRelay(resp any) (*model.LLMResponseSchema, error) {
	switch resp.(type) {
	case openai.ChatCompletionResponse:
		return chatGPTTranslateToRelay(resp.(openai.ChatCompletionResponse))
	default:
		return nil, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
func chatGPTTranslateToRelay(resp openai.ChatCompletionResponse) (*model.LLMResponseSchema, error) {
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
