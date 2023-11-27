package respTransformer

import (
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
	"strconv"
)

func ChatGLMToGrpcCompletion(res ChatGLM.ChatCompletionResponse) (*model.LLMResponseSchema, error) {
	content, err := strconv.Unquote(res.Data.Choices[0].Content)
	if err != nil {
		content = res.Data.Choices[0].Content
	}
	retMessage := model.LLMChatCompletionMessage{
		Role:    res.Data.Choices[0].Role,
		Content: content,
	}
	count := model.LLMTokenCount{
		TotalToken:      int64(res.Data.Usage.TotalTokens),
		PromptToken:     int64(res.Data.Usage.PromptTokens),
		CompletionToken: int64(res.Data.Usage.CompletionTokens),
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}

func ChatGLMToGrpcStream(llm chan string, result chan ChatGLM.ChatCompletionStreamFinishResponse, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	for {
		select {
		case chunk := <-llm:
			resp := model.PartialLLMResponse{
				Response: &model.PartialLLMResponse_Content{
					Content: chunk,
				},
			}
			sv.Send(&resp)
		case res := <-result:
			tokenCount := res.Usage.ToGrpc()
			resp := model.PartialLLMResponse{
				Response:      &model.PartialLLMResponse_Done{},
				LlmTokenCount: &tokenCount,
			}
			return sv.Send(&resp)
		}
	}
}
