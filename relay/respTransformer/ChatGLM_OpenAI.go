package respTransformer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/utils"
	"io"
	"strconv"
	"time"
)

func ChatGLMToOpenAICompletion(res ChatGLM.ChatCompletionResponse) openai.ChatCompletionResponse {
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
	}
}

func ChatGLMToOpenAIStream(c *gin.Context, _r *ChatGLM.ChatCompletionStreamResponse) {
	llm, result := _r.LLM, _r.Finish
	utils.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-llm:
			response := openai.ChatCompletionStreamResponse{
				Object:  "chat.completion.chunk",
				Model:   "chatglm",
				Created: time.Now().Unix(),
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
