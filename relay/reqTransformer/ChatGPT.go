package reqTransformer

import (
	"log"

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
	if len(rs.LlmRequestInfo.Functions) == 0 {
		return openai.ChatCompletionRequest{
			Model:       info.GetModel(),
			Messages:    openaiMessages,
			TopP:        float32(info.GetTopP()),
			Temperature: float32(info.GetTemperature()),
			Stream:      true,
		}
	}

	toolChoice := "none"
	if info.UseFunctionCalling {
		toolChoice = "auto"
	}

	tools := make([]openai.Tool, len(info.Functions))
	for i, f := range info.Functions {
		tools[i] = openai.Tool{
			Type: "function",
			Function: openai.FunctionDefinition{
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
		log.Printf("%#v\n", tools[i])
	}
	return openai.ChatCompletionRequest{
		Model:       info.GetModel(),
		Messages:    openaiMessages,
		TopP:        float32(info.GetTopP()),
		Temperature: float32(info.GetTemperature()),
		Stream:      true,
		Tools:       tools,
		ToolChoice:  toolChoice,
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
