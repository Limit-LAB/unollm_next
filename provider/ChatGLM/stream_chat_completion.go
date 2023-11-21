package ChatGLM

import (
	"bytes"
	"encoding/json"
	"limit.dev/unollm/model/zhipu"
	"limit.dev/unollm/utils"
	"net/http"
	"strings"
)

func (c *Client) ChatCompletionStreamingRequest(body zhipu.ChatCompletionRequest, modelName string) (chan string, chan zhipu.ChatCompletionStreamFinishResponse, error) {
	token, err := utils.CreateJWTToken(c.token, jwtExpire)
	if err != nil {
		return nil, nil, err
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", c.base+modelName+"/sse-invoke", bytes.NewReader(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, nil, err
	}

	reader := utils.NewEventStreamReader(resp.Body, 4096)

	llmCh := make(chan string)
	resultCh := make(chan zhipu.ChatCompletionStreamFinishResponse, 1)

	go func() {
		defer resp.Body.Close()
		for reader.Scanner.Scan() {
			kv := strings.Split(reader.Scanner.Text(), "\n")
			switch kv[0] {
			case "event:add":
				llmCh <- kv[2][5:]
			case "event:finish":
				var usage zhipu.ChatCompletionStreamFinishResponse
				json.NewDecoder(strings.NewReader(kv[3][5:])).Decode(&usage)
				resultCh <- usage
			}
		}
	}()

	return llmCh, resultCh, nil
}
