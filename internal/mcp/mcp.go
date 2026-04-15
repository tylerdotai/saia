package mcp

import "fmt"

type Client struct {
	// TODO: MCP client connections
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(name, url string) error {
	// TODO: Connect to MCP server
	return fmt.Errorf("not implemented")
}

func (c *Client) CallTool(server, tool string, args map[string]interface{}) (string, error) {
	// TODO: Call tool on MCP server
	return "", fmt.Errorf("not implemented")
}
