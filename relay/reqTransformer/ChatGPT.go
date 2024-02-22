package reqTransformer

import (
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
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
		// model: chatgpt::{model}
		Model:       info.GetModel()[9:],
		Messages:    openaiMessages,
		TopP:        float32(info.GetTopP()),
		Temperature: float32(info.GetTemperature()),
		Stream:      true,
	}
}

func ChatGPTGrpcCompletionReq(rs *model.LLMRequestSchema) openai.CompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	prompt := ""
	if len(messages) > 0 {
		prompt = messages[len(messages)-1].GetContent()
	}
	return openai.CompletionRequest{
		Model:       info.GetModel(),
		Prompt:      prompt,
		TopP:        float32(info.GetTopP()),
		Temperature: float32(info.GetTemperature()),
		Stream:      true,
	}
}
