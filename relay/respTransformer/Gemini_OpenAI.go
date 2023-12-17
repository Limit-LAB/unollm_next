package respTransformer

import (
	"encoding/json"
	"errors"
	"github.com/Limit-LAB/go-gemini"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/utils"
	"io"
	"time"
)

func GeminiToOpenAIChatCompletionStream(c *gin.Context, s *gemini.GenerateContentStreamer) {
	utils.SetEventStreamHeaders(c)
	defer s.Close()

	c.Stream(func(w io.Writer) bool {
		rsp, err := s.Receive()
		if errors.Is(err, io.EOF) {
			c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
			return false
		}

		if err != nil {
			// c.Render(-1, utils.CustomEvent{Data: "data: [ERROR]"})
			return false
		}
		var choices []openai.ChatCompletionStreamChoice
		for _, content := range rsp.Candidates {
			role, msg, _ := getGeminiContentFromCandidate(content)
			choices = append(choices, openai.ChatCompletionStreamChoice{
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Role:    role,
					Content: msg,
				},
			})
		}
		response := openai.ChatCompletionStreamResponse{
			Object:  "chat.completion.chunk",
			Model:   "gemini",
			Created: time.Now().Unix(),
			Choices: choices,
		}
		jRsp, _ := json.Marshal(response)
		c.Render(-1, utils.CustomEvent{Data: "data: " + string(jRsp)})
		return true
	})
}
