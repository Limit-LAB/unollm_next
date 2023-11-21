package relay

import (
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model/zhipu"
	"strconv"
)

func ChatGLMTranslateToOpenAI(resp any) (openai.ChatCompletionResponse, error) {
	switch resp.(type) {
	case zhipu.ChatCompletionResponse:
		return chatGLMTranslateToOpenAI(resp.(zhipu.ChatCompletionResponse))
	default:
		return openai.ChatCompletionResponse{}, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
func chatGLMTranslateToOpenAI(res zhipu.ChatCompletionResponse) (openai.ChatCompletionResponse, error) {
	content, err := strconv.Unquote(res.Data.Choices[0].Content)
	if err != nil {
		content = res.Data.Choices[0].Content
	}
	return openai.ChatCompletionResponse{
		ID: res.Data.TaskId,
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    res.Data.Choices[0].Role,
					Content: content,
				},
			},
		},
		Usage: openai.Usage{
			PromptTokens:     res.Data.Usage.PromptTokens,
			TotalTokens:      res.Data.Usage.TotalTokens,
			CompletionTokens: res.Data.Usage.CompletionTokens,
		},
	}, nil
}
