package grpcServer

import (
	"fmt"
	"github.com/Limit-LAB/go-gemini"
	"github.com/Limit-LAB/go-gemini/models"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/relay/reqTransformer"
)

func GeminiChatCompletion(cli *gemini.Client, rs *model.LLMRequestSchema) (*model.LLMResponseSchema, error) {
	req := reqTransformer.GeminiGrpcChatCompletionReq(rs)
	generatedContent, err := cli.GenerateContent(models.GeminiModel(rs.GetLlmRequestInfo().Model), &req)
	if err != nil {
		return nil, fmt.Errorf("gemini generate content: %w", err)
	}
	message := ""
	for _, item := range generatedContent.Candidates[0].Content.Parts {
		if item.Text != nil {
			message += *item.Text
		}
	}
	promptToken, err := geminiTokenCount(cli, req.Contents)
	if err != nil {
		return nil, fmt.Errorf("gemini token count, prompt token: %w", err)
	}
	completionToken, err := geminiTokenCount(cli, []models.Content{generatedContent.Candidates[0].Content})
	if err != nil {
		return nil, fmt.Errorf("gemini token count, completion token: %w", err)
	}

	return &model.LLMResponseSchema{
		Message: &model.LLMChatCompletionMessage{
			Role:    string(generatedContent.Candidates[0].Content.Role),
			Content: message,
		},
		LlmTokenCount: &model.LLMTokenCount{
			PromptToken:     int64(*promptToken),
			CompletionToken: int64(*completionToken),
			TotalToken:      int64(*promptToken + *completionToken),
		},
	}, nil
}

func geminiTokenCount(cli *gemini.Client, contents []models.Content) (*int, error) {
	rst, err := cli.CountToken(models.GeminiPro, models.CountTokenRequest{
		Contents: contents,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini count token: %w", err)
	}
	return &rst.TotalTokens, nil
}
