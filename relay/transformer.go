package relay

import (
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/model"
)

type ResponseGrpcTransformer func(any) (*model.LLMResponseSchema, error)
type ResponseOpenAICompletionTransformer func(any) (openai.ChatCompletionResponse, error)
