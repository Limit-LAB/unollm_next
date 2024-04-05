package grpcServer

import (
	"go.limit.dev/unollm/model"
	anthropic "go.limit.dev/unollm/provider/Anthropic"
)

func NewAnthropicClient(info *model.LLMRequestInfo) *anthropic.Client {
	return anthropic.NewClient(info.GetToken())
}

func AnthropicChatCompletion(cli *anthropic.Client, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	req := anthropic.ChatCompletionRequest{
		Model:       rs.LlmRequestInfo.Model,
		Messages:    make([]anthropic.Message, 0),
		MaxTokens:   4096,
		Temperature: rs.LlmRequestInfo.Temperature,
		TopP:        rs.LlmRequestInfo.TopP,
		TopK:        rs.LlmRequestInfo.TopK,
	}
	for _, m := range rs.Messages {
		req.Messages = append(req.Messages, anthropic.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	respp, err := anthropic.ChatCompletion(cli, req)
	if err != nil {
		return nil, err
	}
	res := &model.LLMResponseSchema{
		Message: &model.LLMChatCompletionMessage{
			Role:    respp.Role,
			Content: respp.Content[0].Text,
		},
		LlmTokenCount: &model.LLMTokenCount{
			PromptToken:     int64(respp.Usage.InputTokens),
			CompletionToken: int64(respp.Usage.OutputTokens),
			TotalToken:      int64(respp.Usage.InputTokens + respp.Usage.OutputTokens),
		},
	}

	return res, nil
}

func AnthropicChatCompletionStreaming(cli *anthropic.Client, rs *model.LLMRequestSchema, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	req := anthropic.ChatCompletionRequest{
		Model:     rs.LlmRequestInfo.Model,
		Messages:  make([]anthropic.Message, 0),
		MaxTokens: 4096,
	}
	for _, m := range rs.Messages {
		req.Messages = append(req.Messages, anthropic.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	respp, err := cli.ChatCompletionStreamingRequest(&req)
	if err != nil {
		return err
	}
	go func() {
		for {
			r := <-respp
			if r.Type == "STOP" {
				sv.Send(&model.PartialLLMResponse{
					Response: &model.PartialLLMResponse_Done{},
					LlmTokenCount: &model.LLMTokenCount{
						PromptToken:     int64(r.Delta.Usage.InputTokens),
						CompletionToken: int64(r.Delta.Usage.OutputTokens),
						TotalToken:      int64(r.Delta.Usage.InputTokens + r.Delta.Usage.OutputTokens),
					},
				})
				return
			} else {
				sv.Send(&model.PartialLLMResponse{
					Response: &model.PartialLLMResponse_Content{
						Content: r.Delta.Text,
					},
				})
			}
		}
	}()
	return nil
}
