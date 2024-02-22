package ChatGLM

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"go.limit.dev/unollm/utils"
)

func (c *Client) EmbeddingRequest(body EmbeddingRequest) (result EmbeddingResponse, err error) {

	token, err := utils.CreateJWTToken(c.apiKey, jwtExpire)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/embeddings", bytes.NewReader(reqBody))
	if err != nil {
		return EmbeddingResponse{}, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

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
