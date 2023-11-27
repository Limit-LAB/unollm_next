package grpcServer

import (
	"fmt"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/relay"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ChatGLMChatCompletion(cli *ChatGLM.Client, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	fmt.Println("CHATGLM_LLM_API")

	req := reqTransformer.ChatGLMGrpcChatCompletionReq(rs)

	res, err := relay.ChatGLMChatCompletionRequest(cli, req) // , info.GetModel()

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return respTransformer.ChatGLMToGrpcCompletion(res)
}

func ChatGLMChatCompletionStreaming(cli *ChatGLM.Client, rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	fmt.Println("CHATGLM_LLM_API")

	req := reqTransformer.ChatGLMGrpcChatCompletionReq(rs)
	req.Incremental = true

	res, err := relay.ChatGLMChatCompletionStreamingRequest(cli, req) // , info.GetModel()
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	return respTransformer.ChatGLMToGrpcStream(res, sv)
}
