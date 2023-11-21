package ChatGLM

import (
	"bytes"
	"encoding/json"
	"limit.dev/unollm/utils"
	"net/http"
	"time"
)

const jwtExpire = time.Duration(10000) * time.Second

func (c *Client) ChatCompletion(body ChatCompletionRequest, modelName string) (result ChatCompletionResponse, err error) {
	token, err := utils.CreateJWTToken(c.apiKey, jwtExpire)
	if err != nil {
		return
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", c.base+modelName+"/invoke", bytes.NewReader(reqBody))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	return
}
