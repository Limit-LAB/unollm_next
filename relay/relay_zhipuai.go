package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model"
	"limit.dev/unollm/model/zhipu"
	"limit.dev/unollm/utils"
)

func ChatGLMBlockingRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	messages := rs.GetMessages()
	req := zhipu.ChatCompletionRequest{
		Temperature: float32(info.GetTemperature()),
		TopP:        float32(info.GetTopP()),
	}
	for _, m := range messages {
		if m.GetRole() == "system" {
			req.Prompt = append(req.Prompt, zhipu.ChatCompletionMessage{
				Role:    zhipu.ChatMessageRoleUser,
				Content: m.GetContent(),
			})
			req.Prompt = append(req.Prompt, zhipu.ChatCompletionMessage{
				Role:    zhipu.ChatMessageRoleAssistant,
				Content: "好的，我明白了。",
			})
			continue
		}
		req.Prompt = append(req.Prompt, zhipu.ChatCompletionMessage{
			Role:    m.GetRole(),
			Content: m.GetContent(),
		})
	}

	res, err := utils.GLMBlockingRequest(req, info.GetModel(), info.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !res.Success {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %f, Error msg: %s", res.ErrorCode, res.ErrorMsg))
	}
	
	return chatGLMTranslateToRelay(res)
}
