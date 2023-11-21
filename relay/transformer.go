package relay

import (
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/model"
)

type GrpcTransformer func(any) (*model.LLMResponseSchema, error)
type OpenAICompletionTransformer func(any) (openai.ChatCompletionResponse, error)
