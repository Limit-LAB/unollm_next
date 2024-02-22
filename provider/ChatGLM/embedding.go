package ChatGLM

import (
	"encoding/json"
	"log"
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	return result, nil
}
