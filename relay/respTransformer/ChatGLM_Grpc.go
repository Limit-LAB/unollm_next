package respTransformer

import (
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func ChatGLMToGrpcCompletion(res ChatGLM.ChatCompletionResponse) (*model.LLMResponseSchema, error) {
	// content, err := strconv.Unquote(res.Choices[0].Message.Content)
	content := res.Choices[0].Message.Content
	retMessage := model.LLMChatCompletionMessage{
		Role:    res.Choices[0].Message.Role,
		Content: content,
	}
	count := model.LLMTokenCount{
		TotalToken:      int64(res.Usage.TotalTokens),
		PromptToken:     int64(res.Usage.PromptTokens),
		CompletionToken: int64(res.Usage.CompletionTokens),
	}

	toolCalls := make([]*model.ToolCall, len(res.Choices[0].Message.ToolCalls))
	for i, _ := range res.Choices[0].Message.ToolCalls {
		toolcall := model.ToolCall{
			Id:        res.Choices[0].Message.ToolCalls[i].Id,
			Name:      res.Choices[0].Message.ToolCalls[i].Function.Name,
			Arguments: res.Choices[0].Message.ToolCalls[i].Function.Arguments,
		}
		toolCalls[i] = &toolcall
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		ToolCalls:     toolCalls,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}

func ChatGLMToGrpcStream(_r *ChatGLM.ChatCompletionStreamingResponse, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	llm, result := _r.ResponseCh, _r.FinishCh
	for {
		select {
		case chunk := <-llm:
			toolCalls := make([]*model.ToolCall, len(chunk.Choices[0].Delta.ToolCalls))
			for i, _ := range chunk.Choices[0].Delta.ToolCalls {
				toolcall := model.ToolCall{
					Id:        chunk.Choices[0].Delta.ToolCalls[i].Id,
					Name:      chunk.Choices[0].Delta.ToolCalls[i].Function.Name,
					Arguments: chunk.Choices[0].Delta.ToolCalls[i].Function.Arguments,
				}
				toolCalls[i] = &toolcall
			}

			resp := model.PartialLLMResponse{
				ToolCalls: toolCalls,
				Response: &model.PartialLLMResponse_Content{
					Content: chunk.Choices[0].Delta.Content,
				},
			}
			sv.Send(&resp)
		case res := <-result:
			tokenCount := res
			resp := model.PartialLLMResponse{
				Response: &model.PartialLLMResponse_Done{},
				LlmTokenCount: &model.LLMTokenCount{
					TotalToken:      int64(tokenCount.TotalTokens),
					PromptToken:     int64(tokenCount.PromptTokens),
					CompletionToken: int64(tokenCount.CompletionTokens),
				},
			}
			return sv.Send(&resp)
		}
	}
}

func ChatGLMToGrpcEmbedding(req *model.EmbeddingRequest, res ChatGLM.EmbeddingResponse) (*model.EmbeddingResponse, error) {
	return &model.EmbeddingResponse{
		EmbeddingRequestInfo: req.GetEmbeddingRequestInfo(),
		Dimension:            1024,
		Vectors:              res.Data[0].Embedding,
	}, nil
}
