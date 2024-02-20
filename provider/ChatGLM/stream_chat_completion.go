package ChatGLM

import (
	"encoding/json"
	"log"
	"strings"

	"go.limit.dev/unollm/utils"
)

type ChatCompletionStreamingResponse struct {
	ResponseCh chan ChatCompletionStreamResponse
	FinishCh   chan Usage
}

func (c *Client) ChatCompletionStreamingRequest(body ChatCompletionRequest) (*ChatCompletionStreamingResponse, error) {
	body.Stream = true

	req, err := c.createRequest(body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	reader := utils.NewEventStreamReader(resp.Body, 4096)

	respCh := make(chan ChatCompletionStreamResponse, c.RespBuf)
	finishCh := make(chan Usage, 1)

	go func() {
		defer resp.Body.Close()
		for reader.Scanner.Scan() {
			kv := strings.Split(reader.Scanner.Text(), "\n")
			if kv[0] == "data: [DONE]" {
				break
			}
			if kv[0][0:6] != "data: " {
				log.Println(kv[0])
				finishCh <- Usage{}
				return
			}
			jsonString := kv[0][6:]
			var result ChatCompletionStreamResponse
			err = json.Unmarshal([]byte(jsonString), &result)
			if err != nil {
				log.Println(err)
				finishCh <- Usage{}
				return
			}
			respCh <- result
			if result.Choices[0].FinishReason != FinishReasonNone {
				finishCh <- result.Usage
			}
		}
	}()

	return &ChatCompletionStreamingResponse{
		ResponseCh: respCh,
		FinishCh:   finishCh,
	}, nil
}
