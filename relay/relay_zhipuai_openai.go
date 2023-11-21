package relay

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"limit.dev/unollm/model/zhipu"
	"limit.dev/unollm/utils"
	"strconv"
)

func ChatGLM2OpenAI(resp any) (openai.ChatCompletionResponse, error) {
	switch resp.(type) {
	case zhipu.ChatCompletionResponse:
		return chatGlm2OpenAI(resp.(zhipu.ChatCompletionResponse))
	default:
		return openai.ChatCompletionResponse{}, status.Errorf(codes.Internal, "ChatGPTTranslateToRelay: resp type is not openai.ChatCompletionResponse")
	}
}
func chatGlm2OpenAI(res zhipu.ChatCompletionResponse) (openai.ChatCompletionResponse, error) {
	content, err := strconv.Unquote(res.Data.Choices[0].Content)
	if err != nil {
		content = res.Data.Choices[0].Content
	}
	return openai.ChatCompletionResponse{
		ID: res.Data.TaskId,
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    res.Data.Choices[0].Role,
					Content: content,
				},
			},
		},
		Usage: openai.Usage{
			PromptTokens:     res.Data.Usage.PromptTokens,
			TotalTokens:      res.Data.Usage.TotalTokens,
			CompletionTokens: res.Data.Usage.CompletionTokens,
		},
	}, nil
}

func chatGlmStream2OpenAI(c *gin.Context, llm chan string, result chan zhipu.ChatCompletionStreamFinishResponse) {
	utils.SetEventStreamHeaders(c)
	// TODO: Stop chan?
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-llm:
			response := openai.ChatCompletionStreamResponse{
				Object: "chat.completion.chunk",
				Model:  "chatglm",
				Choices: []openai.ChatCompletionStreamChoice{
					{
						Delta: openai.ChatCompletionStreamChoiceDelta{
							Content: data,
						},
					},
				},
			}
			jsonResponse, _ := json.Marshal(response)
			c.Render(-1, utils.CustomEvent{Data: "data: " + string(jsonResponse)})
			return true
		case _ = <-result:
			c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
			return false
			//response := openai.ChatCompletionStreamResponse{
			//	Model:   "chatglm",
			//	Object:  "chat.completion.chunk",
			//	Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: ""}}},
			//}
			//c.Render(-1, utils.CustomEvent{Data: "data: " + string(jsonResponse)})
			//return true
		}
	})
}
