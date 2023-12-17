package reqTransformer

import (
	"github.com/Limit-LAB/go-gemini/models"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
)

func GeminiGrpcChatCompletionReq(rs *model.LLMRequestSchema) models.GenerateContentRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()
	req := models.NewGenerateContentRequest().WithGenerationConfig(*models.NewGenerationConfig().
		WithTemperature(float32(info.GetTemperature())).
		WithTopP(float32(info.GetTopP())).
		WithTopK(float32(info.GetTopK())),
	)

	for _, m := range messages {
		if m.GetRole() == "system" {
			req.AppendContent(models.RoleUser, models.NewTextPart(m.GetContent()))
			req.AppendContent(models.RoleModel, models.NewTextPart("好的，我明白了。"))
			continue
		}
		if m.GetRole() == "user" {
			req.AppendContent(models.RoleUser, models.NewTextPart(m.GetContent()))
			continue
		}
		req.AppendContent(models.RoleModel, models.NewTextPart(m.GetContent()))
	}
	return *req
}

func GeminiFromOpenAIChatCompletionReq(req openai.ChatCompletionRequest) models.GenerateContentRequest {
	greq := models.NewGenerateContentRequest().WithGenerationConfig(*models.NewGenerationConfig().
		WithTemperature(req.Temperature).
		WithTopP(req.TopP))

	for _, m := range req.Messages {
		if m.Role == openai.ChatMessageRoleSystem {
			greq.AppendContent(models.RoleUser, models.NewTextPart(m.Content))
			greq.AppendContent(models.RoleModel, models.NewTextPart("好的，我明白了。"))
			continue
		}
		if m.Role == openai.ChatMessageRoleUser {
			greq.AppendContent(models.RoleUser, models.NewTextPart(m.Content))
			continue
		}
		if m.Role == openai.ChatMessageRoleAssistant {
			greq.AppendContent(models.RoleModel, models.NewTextPart(m.Content))
			continue
		}
	}
	return *greq
}
