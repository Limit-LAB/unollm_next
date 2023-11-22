package reqTransformer

import (
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/model"
)

func ChatGPTGrpcChatCompletionReq(rs *model.LLMRequestSchema) openai.ChatCompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	var openaiMessages []openai.ChatCompletionMessage
	for _, m := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}
	return openai.ChatCompletionRequest{
		Model:       info.GetModel(),
		Messages:    openaiMessages,
		TopP:        float32(info.GetTopP()),
		Temperature: float32(info.GetTemperature()),
	}
}
