package anthropic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ChatCompletion(cli *Client, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(
		"POST",
		"https://api.anthropic.com/v1/messages",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return nil, err
	}
	httpreq.Header.Set("x-api-key", cli.key)
	httpreq.Header.Set("anthropic-version", "2023-06-01")
	httpreq.Header.Set("content-type", "application/json")
	resp, err := cli.hc.Do(httpreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respp ChatCompletionResponse
	err = json.Unmarshal(bts, &respp)
	if err != nil {
		return nil, err
	}
	return &respp, nil
}
