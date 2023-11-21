package ChatGLM

import (
	"net/http"
)

type Client struct {
	apiKey string
	base   string
	hc     *http.Client
	token  string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey, base: "https://open.bigmodel.cn/api/paas/v3/model-api/", hc: &http.Client{}}
}

func (c *Client) SetBase(base string) {
	c.base = base
}
