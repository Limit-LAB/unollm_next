package reqTransformer

import (
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMGrpcChatCompletionReq(rs *model.LLMRequestSchema) ChatGLM.ChatCompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	req := ChatGLM.ChatCompletionRequest{
		Temperature: float32(info.GetTemperature()),
		TopP:        float32(info.GetTopP()),
	}
	for _, m := range messages {
		if m.GetRole() == "system" {
			req.Prompt = append(req.Prompt, ChatGLM.ChatCompletionMessage{
				Role:    ChatGLM.ChatMessageRoleUser,
				Content: m.GetContent(),
			})
			req.Prompt = append(req.Prompt, ChatGLM.ChatCompletionMessage{
				Role:    ChatGLM.ChatMessageRoleAssistant,
				Content: "好的，我明白了。",
			})
			continue
		}
		req.Prompt = append(req.Prompt, ChatGLM.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}
	return req
}
