package Baichuan

import (
	"net/http"
)

type Client struct {
	base   string
	hc     *http.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		base:   "https://api.baichuan-ai.com/v1",
		hc:     &http.Client{},
	}
}

func (c *Client) SetBase(base string) {
	c.base = base
}
