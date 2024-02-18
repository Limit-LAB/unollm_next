package respTransformer

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/utils"
)

func ChatGLMToOpenAICompletion(res ChatGLM.ChatCompletionResponse) openai.ChatCompletionResponse {
	content := res.Choices[0].Message.Content
	return openai.ChatCompletionResponse{
		ID: res.Id,
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    res.Choices[0].Message.Role,
					Content: content,
				},
			},
		},
		Usage: openai.Usage{
			PromptTokens:     res.Usage.PromptTokens,
			TotalTokens:      res.Usage.TotalTokens,
			CompletionTokens: res.Usage.CompletionTokens,
		},
	}
}

func ChatGLMToOpenAIStream(c *gin.Context, _r *ChatGLM.ChatCompletionStreamingResponse) {
	llm, result := _r.ResponseChannle, _r.FinishUsageChannle
	utils.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-llm:
			response := openai.ChatCompletionStreamResponse{
				Object:  "chat.completion.chunk",
				Model:   data.Model,
				Created: time.Now().Unix(),
				Choices: []openai.ChatCompletionStreamChoice{
					{
						Delta: openai.ChatCompletionStreamChoiceDelta{
							Content: data.Choices[0].Delta.Content,
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
		}
	})
}
