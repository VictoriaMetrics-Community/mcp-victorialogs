package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	serverMode  string
	sseAddr     string
	entrypoint  string
	bearerToken string

	entryPointURL *url.URL
}

func InitConfig() (*Config, error) {
	result := &Config{
		serverMode:  os.Getenv("MCP_SERVER_MODE"),
		sseAddr:     os.Getenv("MCP_SSE_ADDR"),
		entrypoint:  os.Getenv("VL_INSTANCE_ENTRYPOINT"),
		bearerToken: os.Getenv("VL_INSTANCE_BEARER_TOKEN"),
	}
	if result.entrypoint == "" {
		return nil, fmt.Errorf("VL_INSTANCE_ENTRYPOINT is not set")
	}
	if result.serverMode != "" && result.serverMode != "stdio" && result.serverMode != "sse" {
		return nil, fmt.Errorf("MCP_SERVER_MODE must be 'stdio' or 'sse'")
	}
	if result.serverMode == "" {
		result.serverMode = "stdio"
	}
	if result.sseAddr == "" {
		result.sseAddr = "localhost:8081"
	}

	var err error
	result.entryPointURL, err = url.Parse(result.entrypoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL from VL_INSTANCE_ENTRYPOINT: %w", err)
	}

	return result, nil
}

func (c *Config) IsStdio() bool {
	return c.serverMode == "stdio"
}

func (c *Config) IsSSE() bool {
	return c.serverMode == "sse"
}

func (c *Config) SSEAddr() string {
	return c.sseAddr
}

func (c *Config) BearerToken() string {
	return c.bearerToken
}

func (c *Config) EntryPointURL() *url.URL {
	return c.entryPointURL
}

func (c *Config) AdminAPIURL(path ...string) string {
	return c.entryPointURL.JoinPath(path...).String()
}

func (c *Config) SelectAPIURL(path ...string) string {
	return c.entryPointURL.JoinPath("select", "logsql").JoinPath(path...).String()
}
