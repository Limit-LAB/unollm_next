package relay

import (
	"context"
	"errors"
	"fmt"
	"io"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"

	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.limit.dev/unollm/model"
)

// TODO: read max_tokens, n, stop, frequency_penalty, presence_penalty from meta

func OpenAIChatCompletionRequest(ctx context.Context, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	info := rs.GetLlmRequestInfo()
	fmt.Println("OPENAI_LLM_API")
	config := openai.DefaultConfig(info.GetToken())
	config.BaseURL = info.GetUrl()
	client := openai.NewClientWithConfig(config)

	req := reqTransformer.ChatGPTGrpcChatCompletionReq(rs)

	resp, err := client.CreateChatCompletion(
		ctx,
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

	go func() {
		defer resp.Close()
		i := 0
		for {
			response, err := resp.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("\nStream finished")
				sv.Send(&model.PartialLLMResponse{
					Response: &model.PartialLLMResponse_Done{},
					LlmTokenCount: &model.LLMTokenCount{
						CompletionToken: int64(i),
					},
				})
				return
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}

			if len(response.Choices) != 0 {
				message := response.Choices[0].Delta.Content
				fmt.Printf("\nStream message %d: %s\n", i, message)
				i++
				pr := model.PartialLLMResponse{
					Response: &model.PartialLLMResponse_Content{
						Content: message,
					},
				}
				sv.Send(&pr)
			}
		}
	}()
	return nil
}
