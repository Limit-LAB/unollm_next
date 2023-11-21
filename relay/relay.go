package relay

import (
	context "context"
	"fmt"
	"limit.dev/unollm/model/unoLlmMod"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type UnoForwardServer struct {
}

const OPENAI_LLM_API = "openai"
const CHATGLM_LLM_API = "chatglm"
const AZURE_OPENAI_LLM_API = "azure_openai"
const BAICHUAN_LLM_API = "baichuan"

func (UnoForwardServer) BlockingRequestLLM(ctx context.Context, rs *unoLlmMod.LLMRequestSchema) (*unoLlmMod.LLMResponseSchema, error) {
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

func (UnoForwardServer) StreamRequestLLM(rs *unoLlmMod.LLMRequestSchema, sv unoLlmMod.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	switch info.GetLlmApiType() {
	case CHATGLM_LLM_API:
		return ChatGLMStreamingRequestLLM(rs, sv)
	}
	return status.Errorf(codes.Unimplemented, "method StreamRequestLLM not implemented")
}
func (UnoForwardServer) mustEmbedUnimplementedUnoLLMv1Server() {}
