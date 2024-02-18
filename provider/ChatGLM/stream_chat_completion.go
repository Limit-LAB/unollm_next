package ChatGLM

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.limit.dev/unollm/utils"
)

type ChatCompletionStreamingResponse struct {
	ResponseChannle    chan ChatCompletionStreamResponse
	FinishUsageChannle chan Usage
}

func (c *Client) ChatCompletionStreamingRequest(body ChatCompletionRequest) (*ChatCompletionStreamingResponse, error) {
	body.Stream = true

	token, err := utils.CreateJWTToken(c.apiKey, jwtExpire)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.base, bytes.NewReader(reqBody))
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

	llmCh := make(chan ChatCompletionStreamResponse)
	resultCh := make(chan Usage, 1)

	go func() {
		defer resp.Body.Close()
		for reader.Scanner.Scan() {
			kv := strings.Split(reader.Scanner.Text(), "\n")
			if kv[0] == "data: [DONE]" {
				break
			}
			if kv[0][0:6] != "data: " {
				log.Println(kv[0])
				resultCh <- Usage{}
				return
			}
			json_string := kv[0][6:]
			var result ChatCompletionStreamResponse
			err := json.Unmarshal([]byte(json_string), &result)
			if err != nil {
				log.Println(err)
				resultCh <- Usage{}
				return
			}
			llmCh <- result
			if result.Choices[0].FinishReason != "" {
				resultCh <- result.Usage
			}
		}
	}()

	return &ChatCompletionStreamingResponse{
		ResponseChannle:    llmCh,
		FinishUsageChannle: resultCh,
	}, nil
}
