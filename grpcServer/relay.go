package grpcServer

import (
	"context"
	"log"

	"go.limit.dev/unollm/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnoForwardServer struct {
	model.UnimplementedUnoLLMv1Server
}

type UnoEmbeddingForwardServer struct {
	model.UnimplementedUnoEmbeddingv1Server
}

// mustEmbedUnimplementedUnoLLMv1Server implements model.UnoLLMv1Server.

var _ model.UnoLLMv1Server = (*UnoForwardServer)(nil)
var _ model.UnoEmbeddingv1Server = (*UnoEmbeddingForwardServer)(nil)

const OPENAI_LLM_API = "openai"
const CHATGLM_LLM_API = "chatglm"
const AZURE_OPENAI_LLM_API = "azure_openai"
const BAICHUAN_LLM_API = "baichuan"
const GEMINI_LLM_API = "gemini"

func (uno *UnoForwardServer) BlockingRequestLLM(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	switch info.GetLlmApiType() {
	case OPENAI_LLM_API:
		cli := NewOpenAIClient(info)
		return OpenAIChatCompletion(cli, rs)

	case CHATGLM_LLM_API:
		cli := NewChatGLMClient(info)
		return ChatGLMChatCompletion(cli, rs)

	case GEMINI_LLM_API:
		cli := NewGeminiClient(info)
		return GeminiChatCompletion(cli, rs)

	case AZURE_OPENAI_LLM_API:
		log.Println("AZURE_OPENAI_LLM_API")
		return nil, status.Errorf(codes.Unimplemented, "method BlockingRequestLLM not implemented")

	case BAICHUAN_LLM_API:
		cli := NewBaichuanClient(info)
		if functionCallingRequestMake(rs) {
			res, err := BaichuanChatCompletion(cli, rs)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			functionCallingResponseHandle(res)
			return res, nil
		}
		return BaichuanChatCompletion(cli, rs)
	}

	return nil, status.Errorf(codes.Unimplemented, "LLM for platform %s not implemented", info.GetLlmApiType())
}

func (uno *UnoForwardServer) StreamRequestLLM(rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	switch info.GetLlmApiType() {
	case OPENAI_LLM_API:
		cli := NewOpenAIClient(info)
		return OpenAIChatCompletionStreaming(cli, rs, sv)
	case CHATGLM_LLM_API:
		cli := NewChatGLMClient(info)
		return ChatGLMChatCompletionStreaming(cli, rs, sv)
	case GEMINI_LLM_API:
		cli := NewGeminiClient(info)
		return GeminiChatCompletionStreaming(cli, rs, sv)
	case BAICHUAN_LLM_API:
		cli := NewBaichuanClient(info)
		if functionCallingRequestMake(rs) {
			res, err := BaichuanChatCompletion(cli, rs)
			if err != nil {
				return status.Errorf(codes.Internal, err.Error())
			}
			functionCallingResponseHandle(res)
			functionCallingResponseToStream(res, sv)
			return nil
		}
		return BaichuanChatCompletionStream(cli, rs, sv)
	}
	return status.Errorf(codes.Unimplemented, "method StreamRequestLLM not implemented")
}

func (uno *UnoEmbeddingForwardServer) EmbeddingRequestLLM(ctx context.Context, req *model.EmbeddingRequest) (res *model.EmbeddingResponse, err error) {
	info := req.GetEmbeddingRequestInfo()

	switch info.GetLlmApiType() {
	case "chatglm":
		cli := NewChatGLMClient(&model.LLMRequestInfo{
			LlmApiType: info.GetLlmApiType(),
			Model:      info.GetModel(),
			Url:        info.GetUrl(),
			Token:      info.GetToken(),
		})
		return ChatGLMEmbeddingRequest(cli, req)
	case "openai":
		cli := NewOpenAIClient(&model.LLMRequestInfo{
			LlmApiType: info.GetLlmApiType(),
			Model:      info.GetModel(),
			Url:        info.GetUrl(),
			Token:      info.GetToken(),
		})
		return OpenAIEmbeddingRequest(cli, req)
	}

	return nil, status.Errorf(codes.Unimplemented, "method EmbeddingRequestLLM not implemented")
}
