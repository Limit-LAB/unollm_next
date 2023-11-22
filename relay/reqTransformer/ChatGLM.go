package reqTransformer

import (
	"limit.dev/unollm/model"
	"limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMGrpcChatCompletionReq(rs *model.LLMRequestSchema) ChatGLM.ChatCompletionRequest {
	return ChatGLM.RequestFromLLMRequest(rs)
}
