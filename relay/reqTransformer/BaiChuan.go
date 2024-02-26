package reqTransformer

import (
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/Baichuan"
)

func BaiChuanGrpcChatCompletionReq(rs *model.LLMRequestSchema) Baichuan.ChatCompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	req := Baichuan.ChatCompletionRequest{
		Temperature: float32(info.GetTemperature()),
		TopP:        float32(info.GetTopP()),
		TopK:        int(info.GetTopK()),
	}
	for _, m := range messages {
		if m.GetRole() == "system" {
			req.Messages = append(req.Messages, Baichuan.Message{
				Role:    Baichuan.RoleUser,
				Content: m.GetContent(),
			})
			req.Messages = append(req.Messages, Baichuan.Message{
				Role:    Baichuan.RoleAssistant,
				Content: "好的，我明白了。",
			})
			continue
		}
		req.Messages = append(req.Messages, Baichuan.Message{
			Role:    Baichuan.RoleType(m.GetRole()),
			Content: m.GetContent(),
		})
	}
	return req
}

func BaiChuanFromOpenAIChatCompletionReq(req openai.ChatCompletionRequest) Baichuan.ChatCompletionRequest {
	zpReq := Baichuan.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
	}
	for _, m := range req.Messages {
		if m.Role == openai.ChatMessageRoleSystem {
			zpReq.Messages = append(zpReq.Messages, Baichuan.Message{
				Role:    Baichuan.RoleUser,
				Content: m.Content,
			})
			zpReq.Messages = append(zpReq.Messages, Baichuan.Message{
				Role:    Baichuan.RoleAssistant,
				Content: "好的，我明白了。",
			})
		}
		zpReq.Messages = append(zpReq.Messages, Baichuan.Message{
			Role:    Baichuan.RoleType(m.Role),
			Content: m.Content,
		})
	}
	return zpReq
}
