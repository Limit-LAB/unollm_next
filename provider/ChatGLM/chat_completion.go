package ChatGLM

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"go.limit.dev/unollm/utils"
)

const jwtExpire = time.Duration(10000) * time.Second

var ErrNotSuccess = errors.New("request result not success")

func (c *Client) ChatCompletion(body ChatCompletionRequest) (result ChatCompletionResponse, err error) {
	token, err := utils.CreateJWTToken(c.apiKey, jwtExpire)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	req, err := http.NewRequest("POST", c.base, bytes.NewReader(reqBody))
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	defer resp.Body.Close()
	log.Println("ChatGLM response status: ", resp.Status)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	// if !result.Success {
	// 	err = fmt.Errorf("chatGLM response success is false Error code: %d, Error msg: %s", result.ErrorCode, result.ErrorMsg)
	// 	return ChatCompletionResponse{}, err
	// }

	return result, nil
}
