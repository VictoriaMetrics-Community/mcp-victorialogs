package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	serverMode    string
	listenAddr    string
	entrypoint    string
	bearerToken   string
	disabledTools map[string]bool

	entryPointURL *url.URL
}

func InitConfig() (*Config, error) {
	disabledTools := os.Getenv("MCP_DISABLED_TOOLS")
	disabledToolsMap := make(map[string]bool)
	if disabledTools != "" {
		for _, tool := range strings.Split(disabledTools, ",") {
			tool = strings.Trim(tool, " ,")
			if tool != "" {
				disabledToolsMap[tool] = true
			}
		}
	}
	result := &Config{
		serverMode:    strings.ToLower(os.Getenv("MCP_SERVER_MODE")),
		listenAddr:    os.Getenv("MCP_LISTEN_ADDR"),
		entrypoint:    os.Getenv("VL_INSTANCE_ENTRYPOINT"),
		bearerToken:   os.Getenv("VL_INSTANCE_BEARER_TOKEN"),
		disabledTools: disabledToolsMap,
	}
	// Left for backward compatibility
	if result.listenAddr == "" {
		result.listenAddr = os.Getenv("MCP_SSE_ADDR")
	}
	if result.entrypoint == "" {
		return nil, fmt.Errorf("VL_INSTANCE_ENTRYPOINT is not set")
	}
	if result.serverMode != "" && result.serverMode != "stdio" && result.serverMode != "sse" && result.serverMode != "http" {
		return nil, fmt.Errorf("MCP_SERVER_MODE must be 'stdio', 'sse' or 'http'")
	}
	if result.serverMode == "" {
		result.serverMode = "stdio"
	}
	if result.listenAddr == "" {
		result.listenAddr = "localhost:8081"
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

func (c *Config) ServerMode() string {
	return c.serverMode
}

func (c *Config) ListenAddr() string {
	return c.listenAddr
}

func (c *Config) BearerToken() string {
	return c.bearerToken
}

func (c *Config) EntryPointURL() *url.URL {
	return c.entryPointURL
}

func (c *Config) IsToolDisabled(toolName string) bool {
	if c.disabledTools == nil {
		return false
	}
	disabled, ok := c.disabledTools[toolName]
	return ok && disabled
}
