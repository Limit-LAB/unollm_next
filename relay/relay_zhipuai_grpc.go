package relay

// TODO: n is not supported yet

// TODO: characterglm meta info is not readed from meta

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"limit.dev/unollm/model"
	"strconv"
)

func ChatGLMTranslateToRelay(resp any) (*model.LLMResponseSchema, error) {
	switch resp.(type) {
	case map[string]any:
		return chatGLMTranslateToRelay(resp.(map[string]any))
	default:
		return nil, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
func chatGLMTranslateToRelay(res map[string]any) (*model.LLMResponseSchema, error) {
	data := res["data"].(map[string]interface{})
	choices := data["choices"].([]interface{})
	content, err := strconv.Unquote(choices[0].(map[string]interface{})["content"].(string))
	if err != nil {
		content = choices[0].(map[string]interface{})["content"].(string)
	}
	retMessage := model.LLMChatCompletionMessage{
		Role:    choices[0].(map[string]interface{})["role"].(string),
		Content: content,
	}
	usage := data["usage"].(map[string]interface{})
	count := model.LLMTokenCount{
		TotalToken:      int64(usage["total_tokens"].(float64)),
		PromptToken:     int64(usage["prompt_tokens"].(float64)),
		CompletionToken: int64(usage["completion_tokens"].(float64)),
	}
	retResp := model.LLMResponseSchema{
		Message:       &retMessage,
		LlmTokenCount: &count,
	}
	return &retResp, nil
}
