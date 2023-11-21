package relay

import (
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ChatGLMTranslateToOpenAI(resp any) (openai.ChatCompletionResponse, error) {
	switch resp.(type) {
	case map[string]any:
		return chatGLMTranslateToOpenAI(resp.(map[string]any))
	default:
		return openai.ChatCompletionResponse{}, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
func chatGLMTranslateToOpenAI(res map[string]any) (openai.ChatCompletionResponse, error) {
	panic("not implemented")
}
