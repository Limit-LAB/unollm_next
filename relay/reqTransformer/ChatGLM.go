package reqTransformer

import (
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMGrpcChatCompletionReq(rs *model.LLMRequestSchema) ChatGLM.ChatCompletionRequest {
	return ChatGLM.RequestFromLLMRequest(rs)
}
