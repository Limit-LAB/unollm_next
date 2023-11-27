package ChatGLM

import (
	"bytes"
	"encoding/json"
	"go.limit.dev/unollm/utils"
	"net/http"
	"strings"
)

type ChatCompletionStreamResponse struct {
	LLM    chan string
	Finish chan ChatCompletionStreamFinish
}

func (e *ChatCompletionStreamResponse) Close() {
	if e == nil {
		return
	}
	if e.LLM != nil {
		safeClose[string](e.LLM)
	}
	if e.Finish != nil {
		safeClose[ChatCompletionStreamFinish](e.Finish)
	}
}

func safeClose[T any](ch chan T) {
	defer func() { recover() }()
	if ch != nil {
		close(ch)
	}
}

func (c *Client) ChatCompletionStreamingRequest(body ChatCompletionRequest, modelName string) (*ChatCompletionStreamResponse, error) {
	token, err := utils.CreateJWTToken(c.apiKey, jwtExpire)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.base+modelName+"/sse-invoke", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	reader := utils.NewEventStreamReader(resp.Body, 4096)

	llmCh := make(chan string)
	resultCh := make(chan ChatCompletionStreamFinish, 1)

	go func() {
		defer resp.Body.Close()
		for reader.Scanner.Scan() {
			kv := strings.Split(reader.Scanner.Text(), "\n")
			switch kv[0] {
			case "event:add":
				llmCh <- kv[2][5:]
			case "event:finish":
				var usage ChatCompletionStreamFinish
				json.NewDecoder(strings.NewReader(kv[3][5:])).Decode(&usage)
				resultCh <- usage
			}
		}
	}()

	return &ChatCompletionStreamResponse{LLM: llmCh, Finish: resultCh}, nil
}
