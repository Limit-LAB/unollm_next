package grpcServer

import (
	"log"

	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/Baichuan"
	"go.limit.dev/unollm/relay"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BaichuanChatCompletion(cli *Baichuan.Client, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	log.Println("BAICHUAN_LLM_API")

	req := reqTransformer.BaiChuanGrpcChatCompletionReq(rs)

	res, err := relay.BaiChuanChatCompletionRequest(cli, req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return respTransformer.BaiChuanToGrpcCompletion(res)
}

func BaichuanChatCompletionStream(cli *Baichuan.Client, rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	log.Println("BAICHUAN_LLM_API")

	req := reqTransformer.BaiChuanGrpcChatCompletionReq(rs)

	res, err := relay.BaiChuanChatCompletionStreamingRequest(cli, req)

	if err != nil {
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
	}

	return respTransformer.BaiChuanToGrpcStream(res, sv)
}
