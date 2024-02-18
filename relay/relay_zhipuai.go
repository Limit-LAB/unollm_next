package relay

import (
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMChatCompletionRequest(cli *ChatGLM.Client, req ChatGLM.ChatCompletionRequest) (ChatGLM.ChatCompletionResponse, error) {
	return cli.ChatCompletion(req)
}

func ChatGLMChatCompletionStreamingRequest(cli *ChatGLM.Client, req ChatGLM.ChatCompletionRequest) (*ChatGLM.ChatCompletionStreamingResponse, error) {
	return cli.ChatCompletionStreamingRequest(req)
}
