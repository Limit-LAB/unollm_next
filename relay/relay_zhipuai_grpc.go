package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model/unoLlmMod"
	"limit.dev/unollm/model/zhipu"
	"strconv"
)

func ChatGLM2Grpc(resp any) (*unoLlmMod.LLMResponseSchema, error) {
	switch resp.(type) {
	case zhipu.ChatCompletionResponse:
		return chatGLM2Grpcs(resp.(zhipu.ChatCompletionResponse))
	default:
		return nil, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}

func chatGLM2Grpcs(res zhipu.ChatCompletionResponse) (*unoLlmMod.LLMResponseSchema, error) {
	content, err := strconv.Unquote(res.Data.Choices[0].Content)
	if err != nil {
		content = res.Data.Choices[0].Content
	}
	retMessage := unoLlmMod.LLMChatCompletionMessage{
		Role:    res.Data.Choices[0].Role,
		Content: content,
	}
	count := unoLlmMod.LLMTokenCount{
		TotalToken:      int64(res.Data.Usage.TotalTokens),
		PromptToken:     int64(res.Data.Usage.PromptTokens),
		CompletionToken: int64(res.Data.Usage.CompletionTokens),
	}
	retResp := unoLlmMod.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}

func chatGLMStream2Grpc(llm chan string, result chan zhipu.ChatCompletionStreamFinishResponse, sv unoLlmMod.UnoLLMv1_StreamRequestLLMServer) error {
	for {
		select {
		case chunk := <-llm:
			resp := unoLlmMod.PartialLLMResponse{
				Response: &unoLlmMod.PartialLLMResponse_Content{
					Content: chunk,
				},
			}
			sv.Send(&resp)
		case res := <-result:
			tokenCount := res.Usage.ToGrpc()
			resp := unoLlmMod.PartialLLMResponse{
				Response:      &unoLlmMod.PartialLLMResponse_Done{},
				LlmTokenCount: &tokenCount,
			}
			return sv.Send(&resp)
		}
	}
}
