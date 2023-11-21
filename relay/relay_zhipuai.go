package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"limit.dev/unollm/model/unoLlmMod"
	"limit.dev/unollm/provider/ChatGLM"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model/zhipu"
)

func ChatGLMChatCompletionRequest(ctx context.Context, rs *unoLlmMod.LLMRequestSchema) (*unoLlmMod.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := zhipu.FromLLMRequest(rs)

	cli := ChatGLM.NewClient(info.GetToken())
	res, err := cli.ChatCompletion(req, info.GetModel())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !res.Success {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %d, Error msg: %s", res.ErrorCode, res.ErrorMsg))
	}

	return chatGLM2Grpcs(res)
}

func ChatGLMChatCompletionStreamingRequest(rs *unoLlmMod.LLMRequestSchema, sv unoLlmMod.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := zhipu.FromLLMRequest(rs)
	req.Incremental = true

	cli := ChatGLM.NewClient(info.GetToken())
	llm, result, err := cli.ChatCompletionStreamingRequest(req, info.GetModel())
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return chatGLMStream2Grpc(llm, result, sv)
}
