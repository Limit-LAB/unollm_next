package Baichuan

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var ErrNotSuccess = errors.New("request result not success")

func (c *Client) ChatCompletion(body BaichuanRequestBody) (result BaichuanBlockingResponseBody, err error) {

	reqBody, err := json.Marshal(body)
	if err != nil {
		return BaichuanBlockingResponseBody{}, err

	}

	req, err := http.NewRequest("POST", c.base+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return BaichuanBlockingResponseBody{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.Split(c.apiKey, ".")[0])

	resp, err := c.hc.Do(req)
	if err != nil {
		return BaichuanBlockingResponseBody{}, err

	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return BaichuanBlockingResponseBody{}, err

	}

	return result, err
}
