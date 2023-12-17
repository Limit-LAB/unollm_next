package respTransformer

import (
	"errors"
	"fmt"
	"github.com/Limit-LAB/go-gemini"
	"github.com/Limit-LAB/go-gemini/models"
	"go.limit.dev/unollm/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

func GeminiToGrpcCompletion(resp models.GenerateContentResponse) (*model.LLMResponseSchema, error) {
	role, msg, err := getGeminiContent(resp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Gemini choices is empty")
	}

	retMessage := model.LLMChatCompletionMessage{
		Role:    role,
		Content: msg,
	}
	count := model.LLMTokenCount{}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}

func getGeminiContent(resp models.GenerateContentResponse) (role string, txt string, err error) {
	if len(resp.Candidates) == 0 {
		return "", "", status.Errorf(codes.Internal, "Gemini choices is empty")
	}
	message := resp.Candidates[0]
	return getGeminiContentFromCandidate(message)
}

func getGeminiContentFromCandidate(resp models.GenerateContentCandidate) (role string, txt string, err error) {
	message := resp.Content
	if len(resp.Content.Parts) == 0 {
		return "", "", status.Errorf(codes.Internal, "Gemini choices is empty")
	}
	content := message.Parts[0]
	// FIXME: Gemini Could return Image!
	txt = content.GetText()
	if content.IsInlineData() {
		txt = content.GetInlineData().Data
	}
	if message.Role == models.RoleModel {
		role = "assistant"
	} else {
		role = "user"
	}
	return role, txt, nil
}

func GeminiToGrpcStream(resp *gemini.GenerateContentStreamer, sv model.UnoLLMv1_StreamRequestLLMServer) error {
	defer resp.Close()
	i := 0
	for {
		response, err := resp.Receive()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			sv.Send(&model.PartialLLMResponse{
				Response: &model.PartialLLMResponse_Done{},
				LlmTokenCount: &model.LLMTokenCount{
					CompletionToken: int64(i),
				},
			})
			return nil
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return err
		}

		_, msg, err := getGeminiContent(response)
		pr := model.PartialLLMResponse{
			Response: &model.PartialLLMResponse_Content{
				Content: msg,
			},
		}
		i++
		fmt.Printf("\nStream message %d: %s\n", i, msg)
		sv.Send(&pr)
	}
}
