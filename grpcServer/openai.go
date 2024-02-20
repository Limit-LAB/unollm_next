package grpcServer

import (
	"context"
	"log"

	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/relay"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func OpenAIChatCompletion(cli *openai.Client, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	log.Println("OPENAI_LLM_API")
	req := reqTransformer.ChatGPTGrpcChatCompletionReq(rs)
	resp, err := relay.OpenAIChatCompletionRequest(
		cli,
		req,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return respTransformer.ChatGPTToGrpcCompletion(resp)
}

func OpenAIChatCompletionStreaming(cli *openai.Client, rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	log.Println("OPENAI_LLM_API")

	req := reqTransformer.ChatGPTGrpcChatCompletionReq(rs)
	req.Stream = true

	resp, err := cli.CreateChatCompletionStream(context.Background(), req)

	if err != nil {
		return err
	}

	return respTransformer.ChatGPTToGrpcStream(resp, sv)
}
