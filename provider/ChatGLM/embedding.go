package ChatGLM

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
)

func (c *Client) EmbeddingRequest(body EmbeddingRequest) (result EmbeddingResponse, err error) {
	req, err := c.createRequest(c.url("embeddings"), body)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return EmbeddingResponse{}, err
	}
	defer resp.Body.Close()
	log.Println("ChatGLM response status: ", resp.Status)

	responseBodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmbeddingResponse{}, err
	}
	if resp.StatusCode != 200 {
		slog.Error("unexpected status code: %d", resp.StatusCode, "response body", string(responseBodyString))
		return EmbeddingResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	json.Unmarshal(responseBodyString, &result)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	return result, nil
}
