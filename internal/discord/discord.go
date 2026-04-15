package discord

import (
	"fmt"
)

type Client struct {
	token string
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("Discord token is required")
	}
	return &Client{token: token}, nil
}

func (c *Client) Start() error {
	// TODO: Connect to Discord gateway
	return nil
}

func (c *Client) Stop() error {
	// TODO: Close Discord connection
	return nil
}
