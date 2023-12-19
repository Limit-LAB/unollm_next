package Baichuan

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"go.limit.dev/unollm/utils"
)

func (c *Client) ChatCompletionStreamingRequest(body ChatCompletionRequest) (chan StreamResponse, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err

	}

	req, err := http.NewRequest("POST", c.baseUrl+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+strings.Split(c.apiKey, ".")[0])

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err

	}

	reader := utils.NewEventStreamReader(resp.Body, 4096)

	res := make(chan StreamResponse)

	go func() error {
		defer resp.Body.Close()
		for reader.Scanner.Scan() {
			text := reader.Scanner.Text()
			textJson := text[6:]
			if textJson == "[DONE]" {
				break
			}
			var jjson StreamResponse
			_err := json.NewDecoder(strings.NewReader(textJson)).Decode(&jjson)
			if _err != nil {
				return _err
			}
			// fmt.Println(jjson)
			res <- jjson
		}
		return nil
	}()

	return res, nil
}
