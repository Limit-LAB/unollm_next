package reqTransformer

import (
	"log"

	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
)

func ChatGPTToGrpcRequest(api string, model_type string, token string, req openai.ChatCompletionRequest) (*model.LLMRequestSchema, error) {
	messages := make([]*model.LLMChatCompletionMessage, len(req.Messages))
	for i, _ := range req.Messages {
		messages[i] = &model.LLMChatCompletionMessage{
			Role:    req.Messages[i].Role,
			Content: req.Messages[i].Content,
		}
	}
	tools := make([]*model.Function, len(req.Tools))
	for i, _ := range req.Tools {
		openai_params, ok := req.Tools[i].Function.Parameters.(map[string]any)
		if !ok {
			openai_params = make(map[string]any)
		}
		requireds, ok := openai_params["requireds"].([]string)
		if !ok {
			requireds = make([]string, 0)
		}
		openai_params_properties, ok := openai_params["properties"].(map[string]any)
		if !ok {
			openai_params_properties = make(map[string]any)
		}
		params := make([]*model.FunctionCallingParameter, 0)
		for k, v := range openai_params_properties {
			m, ok := v.(map[string]any)
			if !ok {
				continue
			}
			ty, ok := m["type"].(string)
			param := model.FunctionCallingParameter{}
			param.Name = k
			if ok {
				param.Type = ty
			}
			desc, ok := m["description"].(string)
			if ok {
				param.Description = desc
			}
			enum, ok := m["enum"].([]string)
			if ok {
				param.Enums = enum
			}
			params = append(params, &param)
		}

		tools[i] = &model.Function{
			Name:        req.Tools[i].Function.Name,
			Description: req.Tools[i].Function.Description,
			Parameters:  params,
			Requireds:   requireds,
		}
	}
	usefc := true
	if req.ToolChoice == "none" {
		usefc = false
	}
	return &model.LLMRequestSchema{
		Messages: messages,
		LlmRequestInfo: &model.LLMRequestInfo{
			LlmApiType:         api,
			Model:              model_type,
			Temperature:        float64(req.Temperature),
			TopP:               float64(req.TopP),
			TopK:               float64(0),
			Url:                "TODO: URL",
			Token:              token,
			UseFunctionCalling: usefc,
			Functions:          tools,
		},
	}, nil
}

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
