package Baichuan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var ErrNotSuccess = errors.New("request result not success")

func (c *Client) ChatCompletion(body ChatCompletionRequest) (result ChatCompletionResponse, err error) {

	reqBody, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", c.baseUrl+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.hc.Do(req)
	if err != nil {
		return

	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("Baichuan ChatCompletion: response status code not 200",
			"status_code", resp.StatusCode,
			"response_body", string(responseBody),
		)
		return ChatCompletionResponse{}, ErrNotSuccess
	}
	json.Unmarshal(responseBody, &result)
	if err != nil {
		return ChatCompletionResponse{}, err

	}

	return result, err
}
