package relay

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: read max_tokens, n, stop, frequency_penalty, presence_penalty from meta
func OpenAIChatCompletionRequest(cli *openai.Client, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	req.Stream = false
	return cli.CreateChatCompletion(
		context.Background(),
		req,
	)
}
func OpenAIChatCompletionRequestGrpc(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("OPENAI_LLM_API")
	config := openai.DefaultConfig(info.GetToken())
	config.BaseURL = info.GetUrl()
	client := openai.NewClientWithConfig(config)

	req := reqTransformer.ChatGPTGrpcChatCompletionReq(rs)

	resp, err := OpenAIChatCompletionRequest(
		client,
		req,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return respTransformer.ChatGPTToGrpcCompletion(resp)
}

func OpenAIChatCompletionStreamingRequest(rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	info := rs.GetLlmRequestInfo()
	fmt.Println("OPENAI_LLM_API")
	config := openai.DefaultConfig(info.GetToken())
	config.BaseURL = info.GetUrl()

	client := openai.NewClientWithConfig(config)

	ctx := context.Background()

	req := reqTransformer.ChatGPTGrpcChatCompletionReq(rs)
	req.Stream = true

	resp, err := client.CreateChatCompletionStream(ctx, req)

	if err != nil {
		return err
	}

	return respTransformer.ChatGPTToGrpcStream(resp, sv)
}
