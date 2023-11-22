package relay

import (
	"bufio"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"limit.dev/unollm/utils"
	"net/http"
	"strings"
)

type relayType int

const (
	ChatCompletions relayType = iota
	Completions
)

type CompletionsStreamResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func openaiStreamHandler(c *gin.Context, resp *http.Response, relayMode relayType) string {
	defer resp.Body.Close()
	responseText := ""
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := strings.Index(string(data), "\n"); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	dataChan := make(chan string)
	stopChan := make(chan bool)
	go func() {
		for scanner.Scan() {
			data := scanner.Text()
			if len(data) < 6 { // ignore blank line or wrong format
				continue
			}
			if data[:6] != "data: " && data[:6] != "[DONE]" {
				continue
			}
			dataChan <- data
			data = data[6:]
			if !strings.HasPrefix(data, "[DONE]") {
				switch relayMode {
				case ChatCompletions:
					var streamResponse openai.ChatCompletionStreamResponse
					err := json.Unmarshal([]byte(data), &streamResponse)
					if err != nil {
						continue
					}
					for _, choice := range streamResponse.Choices {
						responseText += choice.Delta.Content
					}
				case Completions:
					var streamResponse CompletionsStreamResponse
					err := json.Unmarshal([]byte(data), &streamResponse)
					if err != nil {
						continue
					}
					for _, choice := range streamResponse.Choices {
						responseText += choice.Text
					}
				}
			}
		}
		stopChan <- true
	}()
	utils.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-dataChan:
			// some implementations may add \r at the end of data
			data = strings.TrimSuffix(data, "\r")
			hasDone := strings.HasPrefix(data, "data: [DONE]")
			if hasDone {
				data = data[:12]
			}
			c.Render(-1, utils.CustomEvent{Data: data})
			return true
		case <-stopChan:
			return false
		}
	})
	return responseText
}
