package anthropic

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.limit.dev/unollm/utils"
)

func (c *Client) ChatCompletionStreamingRequest(body *ChatCompletionRequest) (chan StreamResponse, error) {
	body.Stream = true
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err

	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.key)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("X-Stainless-Stream-Helper", "messages")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err

	}

	reader := utils.NewEventStreamReader(resp.Body, 4096)

	res := make(chan StreamResponse)

	go func() error {
		defer resp.Body.Close()
		var jj map[string]any
		var jjj StreamResponse
		for reader.Scanner.Scan() {
			text := reader.Scanner.Text()
			// read from first \n to the end of the text
			event := text[:strings.Index(text, "\n")]
			text = text[strings.Index(text, "\n")+1+6:]
			log.Println(event)
			log.Println(text)
			if strings.Contains(event, "error") {
				log.Println(text)
				break
			}
			if strings.Contains(event, "message_stop") {
				jjj.Delta.Usage.InputTokens = int(jj["message"].(map[string]any)["usage"].(map[string]any)["input_tokens"].(float64))
				res <- StreamResponse{Type: "STOP", Delta: StreamResponseDelta{
					Usage: jjj.Delta.Usage,
				}}
				return nil
			}
			if jj != nil {
				err = json.NewDecoder(strings.NewReader(text)).Decode(&jjj)
				if err != nil {
					log.Println(err)
					return nil
				}
				res <- jjj
			} else {
				err := json.NewDecoder(strings.NewReader(text)).Decode(&jj)
				if err != nil {
					break
				}
			}
		}
		return reader.Scanner.Err()
	}()
	return res, nil
}
