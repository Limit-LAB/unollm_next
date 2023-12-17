package respTransformer

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/utils"
	"io"
)

// ChatGPTToOpenAIChatCompletionStream ChatComplete
func ChatGPTToOpenAIChatCompletionStream(c *gin.Context, s *openai.ChatCompletionStream) {
	utils.SetEventStreamHeaders(c)
	defer s.Close()

	c.Stream(func(w io.Writer) bool {
		rsp, err := s.Recv()
		if errors.Is(err, io.EOF) {
			c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
			return false
		}
		if err != nil {
			// c.Render(-1, utils.CustomEvent{Data: "data: [ERROR]"})
			return false
		}
		jRsp, _ := json.Marshal(rsp)
		c.Render(-1, utils.CustomEvent{Data: "data: " + string(jRsp)})
		return true
	})
}

func ChatGPTToOpenAICompletionStream(c *gin.Context, s *openai.CompletionStream) {
	utils.SetEventStreamHeaders(c)
	defer s.Close()

	c.Stream(func(w io.Writer) bool {
		rsp, err := s.Recv()
		if errors.Is(err, io.EOF) {
			c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
			return false
		}
		if err != nil {
			// c.Render(-1, utils.CustomEvent{Data: "data: [ERROR]"})
			return false
		}
		jRsp, _ := json.Marshal(rsp)
		c.Render(-1, utils.CustomEvent{Data: "data: " + string(jRsp)})
		return true
	})
}
