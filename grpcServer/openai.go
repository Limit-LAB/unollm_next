package grpcServer

import (
	"context"
	"log"

	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/relay"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"go.limit.dev/unollm/utils"
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
	promptUsage := utils.GetOpenAITokenCount(req.Messages)
	if promptUsage == -1 {
		promptUsage = 0 // TODO: fix token count failed
	}

	resp, err := cli.CreateChatCompletionStream(context.Background(), req)

	if err != nil {
		return err
	}

	return respTransformer.ChatGPTToGrpcStream(promptUsage, resp, sv)
}

type myEmbeddingRequest struct {
	Text string
}

func (emr myEmbeddingRequest) Convert() openai.EmbeddingRequest {

	return openai.EmbeddingRequest{
		Input:          emr.Text,
		Model:          17,
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
	}
}

func OpenAIEmbeddingRequest(cli *openai.Client, req *model.EmbeddingRequest) (*model.EmbeddingResponse, error) {
	log.Println("OPENAI_LLM_API")

	res, err := cli.CreateEmbeddings(context.Background(), myEmbeddingRequest{
		Text: req.Text,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return respTransformer.OpenAIToGrpcEmbedding(res)
}
