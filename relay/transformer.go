package relay

import (
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/model/unoLlmMod"
)

type ResponseGrpcTransformer func(any) (*unoLlmMod.LLMResponseSchema, error)
type ResponseOpenAICompletionTransformer func(any) (openai.ChatCompletionResponse, error)
