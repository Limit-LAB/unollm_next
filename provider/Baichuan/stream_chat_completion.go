package Baichuan

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.limit.dev/unollm/utils"
)

func (c *Client) ChatCompletionStreamingRequest(body BaichuanRequestBody) (chan BaichuanStreamResponseBody, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err

	}

	req, err := http.NewRequest("POST", c.base+"/chat/completions", bytes.NewReader(reqBody))
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

	res := make(chan BaichuanStreamResponseBody)

	go func() {
		defer resp.Body.Close()
		for reader.Scanner.Scan() {
			text := reader.Scanner.Text()
			text_json := text[6:]
			if text_json == "[DONE]" {
				break
			}
			var jjson BaichuanStreamResponseBody
			err := json.NewDecoder(strings.NewReader(text_json)).Decode(&jjson)
			if err != nil {
				panic("fuck")
			}
			log.Println(jjson)
			res <- jjson
		}
	}()

	return res, nil
}
