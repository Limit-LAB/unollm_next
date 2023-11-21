package relay

import (
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model/unoLlmMod"
)

var _ GrpcTransformer = ChatGPT2Grpc

func ChatGPT2Grpc(resp any) (*unoLlmMod.LLMResponseSchema, error) {
	switch resp.(type) {
	case openai.ChatCompletionResponse:
		return chatGPT2Grpc(resp.(openai.ChatCompletionResponse))
	default:
		return nil, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
func chatGPT2Grpc(resp openai.ChatCompletionResponse) (*unoLlmMod.LLMResponseSchema, error) {
	if len(resp.Choices) == 0 {
		return nil, status.Errorf(codes.Internal, "OpenAI choices is empty")
	}
	message := resp.Choices[0].Message
	retMessage := unoLlmMod.LLMChatCompletionMessage{
		Role:    message.Role,
		Content: message.Content,
	}
	count := unoLlmMod.LLMTokenCount{
		TotalToken:      int64(resp.Usage.TotalTokens),
		PromptToken:     int64(resp.Usage.PromptTokens),
		CompletionToken: int64(resp.Usage.CompletionTokens),
	}
	retResp := unoLlmMod.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}
