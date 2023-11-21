package relay

// TODO: system prompt is not supported yet
// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"limit.dev/unollm/model"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/utils"
)

func ChatGLMBlockingRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")
	prompt := make([]interface{}, 0)
	messages := rs.GetMessages()
	for _, m := range messages {
		prompt = append(prompt, map[string]interface{}{
			"role":    m.GetRole(),
			"content": m.GetContent(),
		})
	}

	res, error := utils.GLMBlockingRequest(map[string]interface{}{
		"prompt":      prompt,
		"temperature": info.GetTemperature(),
		"top_p":       info.GetTopP(),
	}, info.GetModel(), info.GetToken())

	if error != nil {
		return nil, status.Errorf(codes.Internal, error.Error())
	} else {
		if !res["success"].(bool) {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %f, Error msg: %s", res["code"], res["msg"]))
		} else {
			data := res["data"].(map[string]interface{})
			choices := data["choices"].([]interface{})
			content, err := strconv.Unquote(choices[0].(map[string]interface{})["content"].(string))
			if err != nil {
				content = choices[0].(map[string]interface{})["content"].(string)
			}
			retMessage := model.LLMChatCompletionMessage{
				Role:    choices[0].(map[string]interface{})["role"].(string),
				Content: content,
			}
			usage := data["usage"].(map[string]interface{})
			count := model.LLMTokenCount{
				TotalToken:      int64(usage["total_tokens"].(float64)),
				PromptToken:     int64(usage["prompt_tokens"].(float64)),
				CompletionToken: int64(usage["completion_tokens"].(float64)),
			}
			retResp := model.LLMResponseSchema{
				Message:       &retMessage,
				LlmTokenCount: &count,
			}
			return &retResp, nil
		}
	}
}
