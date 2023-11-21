package relay

import (
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ (OpenAICompletionTransformer) = ChatGPT2OpenAI

func ChatGPT2OpenAI(resp any) (openai.ChatCompletionResponse, error) {
	switch resp.(type) {
	case openai.ChatCompletionResponse:
		return resp.(openai.ChatCompletionResponse), nil
	default:
		return openai.ChatCompletionResponse{}, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
