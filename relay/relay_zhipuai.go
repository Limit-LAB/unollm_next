package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"limit.dev/unollm/model"
	"limit.dev/unollm/relay/respTransformer"

	"limit.dev/unollm/provider/ChatGLM"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ChatGLMChatCompletionRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := ChatGLM.RequestFromLLMRequest(rs)

	cli := ChatGLM.NewClient(info.GetToken())
	res, err := cli.ChatCompletion(req, info.GetModel())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !res.Success {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %d, Error msg: %s", res.ErrorCode, res.ErrorMsg))
	}

	return respTransformer.ChatGLMToGrpcCompletion(res)
}

func ChatGLMChatCompletionStreamingRequest(rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := ChatGLM.RequestFromLLMRequest(rs)
	req.Incremental = true

	cli := ChatGLM.NewClient(info.GetToken())
	llm, result, err := cli.ChatCompletionStreamingRequest(req, info.GetModel())
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return respTransformer.ChatGLMToGrpcStream(llm, result, sv)
}
