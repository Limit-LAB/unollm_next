package ChatGLM

import (
	"bytes"
	"encoding/json"
	"go.limit.dev/unollm/utils"
	"net/http"
)

type Client struct {
	base    string
	hc      *http.Client
	apiKey  string
	RespBuf int
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		// base:   "https://open.bigmodel.cn/api/paas/v3/model-api/",
		base:    "https://open.bigmodel.cn/api/paas/v4/",
		hc:      &http.Client{},
		RespBuf: 5,
	}
}

func (c *Client) SetBase(base string) {
	c.base = base
}

func (c *Client) createRequest(url string, body any) (*http.Request, error) {
	token, err := utils.CreateJWTToken(c.apiKey, jwtExpire)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
