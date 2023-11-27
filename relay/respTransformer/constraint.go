package respTransformer

import (
	"errors"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func CompletionToOpenAI(resp any) (openai.ChatCompletionResponse, error) {
	switch resp.(type) {
	case ChatGLM.ChatCompletionResponse:
		return ChatGLMToOpenAICompletion(resp.(ChatGLM.ChatCompletionResponse)), nil
	case openai.ChatCompletionResponse:
		return resp.(openai.ChatCompletionResponse), nil
	default:
		return openai.ChatCompletionResponse{}, errors.New("unknown response type")
	}
}

func ChatCompletionGrpc(resp any) (*model.LLMResponseSchema, error) {
	switch resp.(type) {
	case ChatGLM.ChatCompletionResponse:
		return ChatGLMToGrpcCompletion(resp.(ChatGLM.ChatCompletionResponse))
	case openai.ChatCompletionResponse:
		return ChatGPTToGrpcCompletion(resp.(openai.ChatCompletionResponse))
	default:
		return nil, errors.New("unknown response type")
	}
}

func ChatCompletionStreamToGrpc(resp any, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	switch resp.(type) {
	case *ChatGLM.ChatCompletionStreamResponse:
		return ChatGLMToGrpcStream(resp.(*ChatGLM.ChatCompletionStreamResponse), sv)
	case *openai.ChatCompletionStream:
		return ChatGPTToGrpcStream(resp.(*openai.ChatCompletionStream), sv)
	default:
		return errors.New("unknown response type")
	}
}
