package ChatGLM

import (
	"encoding/json"
	"log"
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
	log.Println("ChatGLM response status: ", resp.Status)
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
