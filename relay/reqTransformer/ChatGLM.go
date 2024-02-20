package reqTransformer

import (
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMGrpcChatCompletionReq(rs *model.LLMRequestSchema) ChatGLM.ChatCompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	req := ChatGLM.ChatCompletionRequest{
		Model:       info.GetModel(),
		Temperature: float32(info.GetTemperature()),
		TopP:        float32(info.GetTopP()),
	}
	for _, m := range messages {
		req.Messages = append(req.Messages, ChatGLM.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}
	return req
}

func ChatGLMFromOpenAIChatCompletionReq(req openai.ChatCompletionRequest) ChatGLM.ChatCompletionRequest {
	zpReq := ChatGLM.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
	}

	for _, m := range req.Messages {
		zpReq.Messages = append(zpReq.Messages, ChatGLM.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return zpReq
}
