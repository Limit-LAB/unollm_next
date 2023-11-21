package relay

import (
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/model/unoLlmMod"
)

type GrpcTransformer func(any) (*unoLlmMod.LLMResponseSchema, error)
type OpenAICompletionTransformer func(any) (openai.ChatCompletionResponse, error)
