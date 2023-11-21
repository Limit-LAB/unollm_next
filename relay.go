package unollm

import (
	context "context"
	"fmt"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type UnoForwardServer struct {
}

const OPENAI_LLM_API = "openai"
const CHATGLM_LLM_API = "chatglm"
const AZURE_OPENAI_LLM_API = "azure_openai"
const BAICHUAN_LLM_API = "baichuan"

func (UnoForwardServer) BlockingRequestLLM(ctx context.Context, rs *LLMRequestSchema) (*LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	switch info.GetLlmApiType() {
	case OPENAI_LLM_API:
		return OpenaiBlockingRequest(ctx, rs)
	case CHATGLM_LLM_API:
		return ChatGLMBlockingRequest(ctx, rs)
	case AZURE_OPENAI_LLM_API:
		fmt.Println("AZURE_OPENAI_LLM_API")
		return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")
	case BAICHUAN_LLM_API:
		fmt.Println("BAICHUAN_LLM_API")
		return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")
	}
	return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")
}

func (UnoForwardServer) StreamRequestLLM(rs *LLMRequestSchema, sv UnoLLMv1_StreamRequestLLMServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRequestLLM not implemented")
}
func (UnoForwardServer) mustEmbedUnimplementedUnoLLMv1Server() {}
