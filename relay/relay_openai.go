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

func OpenAICompletionRequest(cli *openai.Client, req openai.CompletionRequest) (openai.CompletionResponse, error) {
	req.Stream = false
	return cli.CreateCompletion(
		context.Background(),
		req,
	)
}

func OpenAIChatCompletionStream(cli *openai.Client, req openai.ChatCompletionRequest) (*openai.ChatCompletionStream, error) {
	req.Stream = true
	return cli.CreateChatCompletionStream(
		context.Background(),
		req,
	)
}

func OpenAICompletionStream(cli *openai.Client, req openai.CompletionRequest) (*openai.CompletionStream, error) {
	req.Stream = true
	return cli.CreateCompletionStream(
		context.Background(),
		req,
	)
}

func OpenAIEmbeddingRequest(cli *openai.Client, req openai.EmbeddingRequest) (openai.EmbeddingResponse, error) {
	return cli.CreateEmbeddings(
		context.Background(),
		req,
	)
}
