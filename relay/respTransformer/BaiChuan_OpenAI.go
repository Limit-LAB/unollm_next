package respTransformer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/Baichuan"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/utils"
	"io"
	"time"
)

func BaiChuanToOpenAICompletion(res Baichuan.ChatCompletionResponse) openai.ChatCompletionResponse {
	return openai.ChatCompletionResponse{
		ID: res.Id,
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    string(res.Choices[0].Message.Role),
					Content: res.Choices[0].Message.Content,
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

func BaiChuanToOpenAIStream(c *gin.Context, _r chan Baichuan.StreamResponse) error {
	defer close(_r)
	utils.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-_r:
			if data.Model == "STOP" {
				c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
				return false
			}
			response := openai.ChatCompletionStreamResponse{
				Object:  "chat.completion.chunk",
				Model:   ChatGLM.ModelGLM3Turbo,
				Created: time.Now().Unix(),
				Choices: []openai.ChatCompletionStreamChoice{
					{
						Delta: openai.ChatCompletionStreamChoiceDelta{
							Content: data.Choices[0].Delta.Content,
							Role:    string(data.Choices[0].Delta.Role),
						},
					},
				},
			}
			jsonResponse, _ := json.Marshal(response)
			c.Render(-1, utils.CustomEvent{Data: "data: " + string(jsonResponse)})
			return true
		}
	})
	return nil
}
