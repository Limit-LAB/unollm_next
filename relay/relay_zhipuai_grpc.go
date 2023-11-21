package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model"
	"limit.dev/unollm/model/zhipu"
	"strconv"
)

func ChatGLMTranslateToRelay(resp any) (*model.LLMResponseSchema, error) {
	switch resp.(type) {
	case zhipu.ChatCompletionResponse:
		return chatGLMTranslateToRelay(resp.(zhipu.ChatCompletionResponse))
	default:
		return nil, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}

func chatGLMTranslateToRelay(res zhipu.ChatCompletionResponse) (*model.LLMResponseSchema, error) {
	content, err := strconv.Unquote(res.Data.Choices[0].Content)
	if err != nil {
		content = res.Data.Choices[0].Content
	}
	retMessage := model.LLMChatCompletionMessage{
		Role:    res.Data.Choices[0].Role,
		Content: content,
	}
	count := model.LLMTokenCount{
		TotalToken:      int64(res.Data.Usage.TotalTokens),
		PromptToken:     int64(res.Data.Usage.PromptTokens),
		CompletionToken: int64(res.Data.Usage.CompletionTokens),
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}
