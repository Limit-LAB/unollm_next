package relay

import (
	"context"
	"fmt"

	"go.limit.dev/unollm/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnoForwardServer struct {
	model.UnimplementedUnoLLMv1Server
}

// mustEmbedUnimplementedUnoLLMv1Server implements model.UnoLLMv1Server.

var _ model.UnoLLMv1Server = (*UnoForwardServer)(nil)

const OPENAI_LLM_API = "openai"
const CHATGLM_LLM_API = "chatglm"
const AZURE_OPENAI_LLM_API = "azure_openai"
const BAICHUAN_LLM_API = "baichuan"

func (uno *UnoForwardServer) BlockingRequestLLM(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	switch info.GetLlmApiType() {
	case OPENAI_LLM_API:
		return OpenAIChatCompletionRequestGrpc(ctx, rs)
	case CHATGLM_LLM_API:
		return ChatGLMChatCompletionRequestGrpc(ctx, rs)
	case AZURE_OPENAI_LLM_API:
		fmt.Println("AZURE_OPENAI_LLM_API")
		return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")
	case BAICHUAN_LLM_API:
		fmt.Println("BAICHUAN_LLM_API")
		return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")
	}
	return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")
}

func (uno *UnoForwardServer) StreamRequestLLM(rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	switch info.GetLlmApiType() {
	case OPENAI_LLM_API:
		return OpenAIChatCompletionStreamingRequest(rs, sv)
	case CHATGLM_LLM_API:
		return ChatGLMChatCompletionStreamingRequestGrpc(rs, sv)
	}
	return status.Errorf(codes.Unimplemented, "method StreamRequestLLM not implemented")
}
