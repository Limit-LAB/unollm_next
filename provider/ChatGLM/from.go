package ChatGLM

import (
	"limit.dev/unollm/model"
)

func RequestFromLLMRequest(rs *model.LLMRequestSchema) ChatCompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	req := ChatCompletionRequest{
		Temperature: float32(info.GetTemperature()),
		TopP:        float32(info.GetTopP()),
	}
	for _, m := range messages {
		if m.GetRole() == "system" {
			req.Prompt = append(req.Prompt, ChatCompletionMessage{
				Role:    ChatMessageRoleUser,
				Content: m.GetContent(),
			})
			req.Prompt = append(req.Prompt, ChatCompletionMessage{
				Role:    ChatMessageRoleAssistant,
				Content: "好的，我明白了。",
			})
			continue
		}
		req.Prompt = append(req.Prompt, ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}
	return req
}
