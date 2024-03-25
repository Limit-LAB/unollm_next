package ChatGLM

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"
)

const jwtExpire = time.Duration(10000) * time.Second

func (c *Client) url(tail string) string {
	return c.base + tail
}

func (c *Client) ChatCompletion(body ChatCompletionRequest) (result ChatCompletionResponse, err error) {
	req, err := c.createRequest(c.url("chat/completions/"), body)
	if err != nil {
		return result, err
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	defer resp.Body.Close()
	responseBodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	if resp.StatusCode != 200 {
		slog.Error("unexpected status code: %d", resp.StatusCode, "response body", string(responseBodyString))
		return ChatCompletionResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Println("ChatGLM response status: ", resp.Status)
	err = json.Unmarshal(responseBodyString, &result)
	return result, err
}
