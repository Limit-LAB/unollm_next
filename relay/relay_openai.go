package relay

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

// TODO: read max_tokens, n, stop, frequency_penalty, presence_penalty from meta
func OpenAIChatCompletionRequest(cli *openai.Client, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	req.Stream = false
	return cli.CreateChatCompletion(
		context.Background(),
		req,
	)
}
