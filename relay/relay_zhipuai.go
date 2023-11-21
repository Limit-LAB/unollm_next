package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"limit.dev/unollm/model/unoLlmMod"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model/zhipu"
	"limit.dev/unollm/utils"
)

func ChatGLMBlockingRequest(ctx context.Context, rs *unoLlmMod.LLMRequestSchema) (*unoLlmMod.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := zhipu.FromLLMRequest(rs)

	res, err := utils.GLMBlockingRequest(req, info.GetModel(), info.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !res.Success {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %d, Error msg: %s", res.ErrorCode, res.ErrorMsg))
	}

	return chatGLMTranslateToRelay(res)
}

func ChatGLMStreamingRequestLLM(rs *unoLlmMod.LLMRequestSchema, sv unoLlmMod.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := zhipu.FromLLMRequest(rs)
	req.Incremental = true

	llm, result, err := utils.GLMStreamingRequest(req, info.GetModel(), info.GetToken())
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	for {
		select {
		case llm_message := <-llm:
			resp := unoLlmMod.PartialLLMResponse{
				Response: &unoLlmMod.PartialLLMResponse_Content{
					Content: llm_message,
				},
			}
			if err = sv.Send(&resp); err != nil {
				return err
			}
		case res := <-result:
			tokenCount := unoLlmMod.LLMTokenCount{
				TotalToken:      int64(res["usage"].(map[string]interface{})["total_tokens"].(float64)),
				PromptToken:     int64(res["usage"].(map[string]interface{})["prompt_tokens"].(float64)),
				CompletionToken: int64(res["usage"].(map[string]interface{})["completion_tokens"].(float64)),
			}
			resp := unoLlmMod.PartialLLMResponse{
				Response:      &unoLlmMod.PartialLLMResponse_Done{},
				LlmTokenCount: &tokenCount,
			}
			if err = sv.Send(&resp); err != nil {
				return err
			}
		}
	}
}
