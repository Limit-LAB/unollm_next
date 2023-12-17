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

func ChatGLMFromOpenAIChatCompletionReq(req openai.ChatCompletionRequest) ChatGLM.ChatCompletionRequest {
	zpReq := ChatGLM.ChatCompletionRequest{
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Incremental: req.Stream,
	}

	for _, m := range req.Messages {
		if m.Role == openai.ChatMessageRoleSystem {
			zpReq.Prompt = append(zpReq.Prompt, ChatGLM.ChatCompletionMessage{
				Role:    ChatGLM.ChatMessageRoleUser,
				Content: m.Content,
			})
			zpReq.Prompt = append(zpReq.Prompt, ChatGLM.ChatCompletionMessage{
				Role:    ChatGLM.ChatMessageRoleAssistant,
				Content: "好的，我明白了。",
			})
		}
		zpReq.Prompt = append(zpReq.Prompt, ChatGLM.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return zpReq
}
