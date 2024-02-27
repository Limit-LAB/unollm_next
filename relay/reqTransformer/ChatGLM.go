package reqTransformer

import (
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMGrpcChatCompletionReq(rs *model.LLMRequestSchema) ChatGLM.ChatCompletionRequest {
	info := rs.GetLlmRequestInfo()
	messages := rs.GetMessages()

	toolChoice := "none"
	if info.UseFunctionCalling {
		toolChoice = "auto"
	}

	tools := make([]ChatGLM.GLMTool, len(info.Functions))
	for i, f := range info.Functions {
		tools[i] = ChatGLM.GLMTool{
			Type: "function",
			Function: ChatGLM.GLMFunction{
				Name:        f.Name,
				Description: f.Description,
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
					"required":   f.Requireds,
				},
			},
		}
		for j := 0; j < len(f.Parameters); j++ {
			tools[i].Function.Parameters.(map[string]any)["properties"].(map[string]any)[f.Parameters[j].Name] = map[string]any{
				"type":        f.Parameters[j].Type,
				"description": f.Parameters[j].Description,
			}
			if len(f.Parameters[j].Enums) != 0 {
				tools[i].Function.Parameters.(map[string]any)["properties"].(map[string]any)[f.Parameters[j].Name].(map[string]any)["enum"] = f.Parameters[j].Enums
			}
		}
	}

	req := ChatGLM.ChatCompletionRequest{
		Model:       info.GetModel(),
		Temperature: float32(info.GetTemperature()),
		TopP:        float32(info.GetTopP()),
		Tools:       tools,
		ToolChoice:  toolChoice,
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
	tools := make([]ChatGLM.GLMTool, len(req.Tools))
	for i, f := range req.Tools {
		functions := ChatGLM.GLMFunction{
			Name:        f.Function.Name,
			Description: f.Function.Description,
			Parameters:  f.Function.Parameters,
		}

		tools[i] = ChatGLM.GLMTool{
			Type:     "function",
			Function: functions,
		}
	}

	toolChoice, ok := req.ToolChoice.(string)
	if !ok {
		toolChoice = "auto"
	}

	zpReq := ChatGLM.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
		ToolChoice:  toolChoice,
		Tools:       tools,
	}

	for _, m := range req.Messages {
		zpReq.Messages = append(zpReq.Messages, ChatGLM.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return zpReq
}

func ChatGLMGrpcEmbeddingReq(req *model.EmbeddingRequest) ChatGLM.EmbeddingRequest {
	return ChatGLM.EmbeddingRequest{
		Input: req.GetText(),
		Model: req.GetEmbeddingRequestInfo().GetModel(),
	}
}
