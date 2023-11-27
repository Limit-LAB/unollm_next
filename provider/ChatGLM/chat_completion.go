package ChatGLM

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.limit.dev/unollm/utils"
	"net/http"
	"time"
)

const jwtExpire = time.Duration(10000) * time.Second

var ErrNotSuccess = errors.New("request result not success")

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
	if err != nil {
		return
	}
	if !result.Success {
		err = fmt.Errorf("chatGLM response success is false Error code: %d, Error msg: %s", result.ErrorCode, result.ErrorMsg)
		return
	}

	return
}
