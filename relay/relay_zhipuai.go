package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMChatCompletionRequest(cli *ChatGLM.Client, req ChatGLM.ChatCompletionRequest) (ChatGLM.ChatCompletionResponse, error) {
	req.Incremental = false
	return cli.ChatCompletion(req, ChatGLM.ModelTurbo) // TODO: read model from meta
}

func ChatGLMChatCompletionStreamingRequest(cli *ChatGLM.Client, req ChatGLM.ChatCompletionRequest) (*ChatGLM.ChatCompletionStreamResponse, error) {
	req.Incremental = true
	return cli.ChatCompletionStreamingRequest(req, ChatGLM.ModelTurbo)
}
