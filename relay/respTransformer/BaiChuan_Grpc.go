package respTransformer

import (
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/Baichuan"
)

func BaiChuanToGrpcCompletion(res Baichuan.ChatCompletionResponse) (*model.LLMResponseSchema, error) {
	retMessage := model.LLMChatCompletionMessage{
		Role:    string(res.Choices[0].Message.Role),
		Content: res.Choices[0].Message.Content,
	}
	count := model.LLMTokenCount{
		TotalToken:      int64(res.Usage.TotalTokens),
		PromptToken:     int64(res.Usage.PromptTokens),
		CompletionToken: int64(res.Usage.CompletionTokens),
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}

func BaiChuanToGrpcStream(_r chan Baichuan.StreamResponse, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	defer close(_r)
	for {
		chunk := <-_r
		if chunk.Model == "STOP" {
			sv.Send(&model.PartialLLMResponse{
				Response: &model.PartialLLMResponse_Done{},
			})
			return nil
		}

		content := model.PartialLLMResponse_Content{
			Content: chunk.Choices[0].Delta.Content,
		}
		sv.Send(&model.PartialLLMResponse{
			Response: &content,
		})
	}
}
