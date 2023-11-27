package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"context"
	"fmt"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"

	"go.limit.dev/unollm/provider/ChatGLM"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ChatGLMChatCompletionRequest(cli *ChatGLM.Client, req ChatGLM.ChatCompletionRequest) (ChatGLM.ChatCompletionResponse, error) {
	res, err := cli.ChatCompletion(req, ChatGLM.ModelTurbo) // TODO: read model from meta
	if err != nil {
		return res, err
	}
	if !res.Success {
		return res, status.Errorf(codes.Internal, fmt.Sprintf("chatGLM response success is false Error code: %d, Error msg: %s", res.ErrorCode, res.ErrorMsg))
	}
	return res, nil
}

func ChatGLMChatCompletionRequestGrpc(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := reqTransformer.ChatGLMGrpcChatCompletionReq(rs)

	cli := ChatGLM.NewClient(info.GetToken())
	res, err := ChatGLMChatCompletionRequest(cli, req) // , info.GetModel()

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return respTransformer.ChatGLMToGrpcCompletion(res)
}

func ChatGLMChatCompletionStreamingRequest(cli *ChatGLM.Client, req ChatGLM.ChatCompletionRequest) (*ChatGLM.ChatCompletionStreamResponse, error) {
	req.Incremental = true
	return cli.ChatCompletionStreamingRequest(req, ChatGLM.ModelTurbo)
}

func ChatGLMChatCompletionStreamingRequestGrpc(rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	fmt.Println("CHATGLM_LLM_API")

	req := reqTransformer.ChatGLMGrpcChatCompletionReq(rs)
	req.Incremental = true

	cli := ChatGLM.NewClient(info.GetToken())
	_r, err := ChatGLMChatCompletionStreamingRequest(cli, req) // , info.GetModel()
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return respTransformer.ChatGLMToGrpcStream(_r, sv)
}
