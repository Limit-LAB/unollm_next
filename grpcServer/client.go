package grpcServer

import (
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func NewOpenAIClient(info *model.LLMRequestInfo) *openai.Client {
	config := openai.DefaultConfig(info.GetToken())
	if info.GetUrl() != "" {
		config.BaseURL = info.GetUrl()
	}

	return openai.NewClientWithConfig(config)
}

func NewChatGLMClient(info *model.LLMRequestInfo) *ChatGLM.Client {
	return ChatGLM.NewClient(info.GetToken())
}
