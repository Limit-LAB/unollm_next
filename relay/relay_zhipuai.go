package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model"
	"limit.dev/unollm/utils"
)

func ChatGLMBlockingRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")
	prompt := make([]interface{}, 0)
	messages := rs.GetMessages()
	for _, m := range messages {
		if m.GetRole() == "system" {
			// TODO: system prompt is not supported yet
			prompt = append(prompt, map[string]interface{}{
				"role":    "user",
				"content": m.GetContent(),
			})
			prompt = append(prompt, map[string]interface{}{
				"role":    "system",
				"content": "好的，我明白了。",
			})
			continue
		}
		prompt = append(prompt, map[string]interface{}{
			"role":    m.GetRole(),
			"content": m.GetContent(),
		})
	}

	res, err := utils.GLMBlockingRequest(map[string]interface{}{
		"prompt":      prompt,
		"temperature": info.GetTemperature(),
		"top_p":       info.GetTopP(),
	}, info.GetModel(), info.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if resSuccess, ok := res["success"].(bool); ok && !resSuccess {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %f, Error msg: %s", res["code"], res["msg"]))
	}
	return chatGLMTranslateToRelay(res)
}
