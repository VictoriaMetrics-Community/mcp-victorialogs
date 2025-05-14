package tools

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/config"
)

func GetTextBodyForRequest(req *http.Request, cfg *config.Config) *mcp.CallToolResult {
	if cfg.BearerToken() != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.BearerToken()))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to do request: %v", err))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %v", err))
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected response status code %v: %s", resp.StatusCode, string(body)))
	}
	return mcp.NewToolResultText(string(body))
}

type ToolReqParamType interface {
	string | float64 | bool | []string
}

func GetToolReqParam[T ToolReqParamType](tcr mcp.CallToolRequest, param string, required bool) (T, error) {
	var value T
	matchArg, ok := tcr.Params.Arguments[param]
	if ok {
		value, ok = matchArg.(T)
		if !ok {
			return value, fmt.Errorf("%s has wrong type: %T", param, matchArg)
		}
	} else if required {
		return value, fmt.Errorf("%s param is required", param)
	}
	return value, nil
}

func GetToolReqTenant(tcr mcp.CallToolRequest) (string, string, error) {
	tenant, err := GetToolReqParam[string](tcr, "tenant", false)
	if err != nil {
		return "", "", fmt.Errorf("failed to get tenant: %v", err)
	}
	tenantParts := strings.Split(tenant, ":")
	if len(tenantParts) > 2 {
		return "", "", fmt.Errorf("tenant must be in the format AccountID:ProjectID")
	}
	accountID := "0"
	projectID := "0"
	if len(tenantParts) > 0 {
		accountID = tenantParts[0]
		if accountID == "" {
			accountID = "0"
		}
	}
	if len(tenantParts) > 1 {
		projectID = tenantParts[1]
		if projectID == "" {
			projectID = "0"
		}
	}
	return accountID, projectID, nil
}
