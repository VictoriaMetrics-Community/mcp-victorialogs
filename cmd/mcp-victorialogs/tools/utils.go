package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/config"
)

func CreateSelectRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	accountID, projectID, err := GetToolReqTenant(tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}

	selectURL, err := getSelectURL(ctx, cfg, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	bearerToken := cfg.BearerToken()
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}
	req.Header.Set("AccountID", accountID)
	req.Header.Set("ProjectID", projectID)

	return req, nil
}

func CreateAdminRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	selectURL, err := getRootURL(ctx, cfg, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	bearerToken := cfg.BearerToken()
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}
	return req, nil
}

func getRootURL(_ context.Context, cfg *config.Config, _ mcp.CallToolRequest, path ...string) (string, error) {
	return cfg.EntryPointURL().JoinPath(path...).String(), nil
}

func getSelectURL(_ context.Context, cfg *config.Config, _ mcp.CallToolRequest, path ...string) (string, error) {
	return cfg.EntryPointURL().JoinPath("select", "logsql").JoinPath(path...).String(), nil
}

func GetTextBodyForRequest(req *http.Request, _ *config.Config) *mcp.CallToolResult {
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
	string | float64 | bool | []string | []any
}

func GetToolReqParam[T ToolReqParamType](tcr mcp.CallToolRequest, param string, required bool) (T, error) {
	var value T
	matchArg, ok := tcr.GetArguments()[param]
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

func ptr[T any](v T) *T {
	return &v
}
