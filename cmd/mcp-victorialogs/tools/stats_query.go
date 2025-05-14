package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/config"
)

var (
	toolStatsQuery = mcp.NewTool("stats_query",
		mcp.WithDescription("Log stats for the given query at the given timestamp (time) in the format compatible with Prometheus querying API. This tool uses `/select/logsql/stats_query` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Querying log stats",
			ReadOnlyHint:    true,
			DestructiveHint: false,
			OpenWorldHint:   true,
		}),
		mcp.WithString("tenant",
			mcp.Title("Tenant name (Account ID and Project ID)"),
			mcp.Description("Name of the tenant for which the data will be displayed (format AccountID:ProjectID)"),
			mcp.DefaultString("0:0"),
			mcp.Pattern(`^([0-9]+)\:[0-9]+$`),
		),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Title("LogsQL expression"),
			mcp.Description(`LogsQL expression that must contain stats pipe. The calculated stats is converted into metrics with labels from by(...) clause of the | stats by(...) pipe.`),
		),
		mcp.WithString("time",
			mcp.Title("Query time"),
			mcp.Description("Time in RFC3339. If it's missing, then it equals to the current time."),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
	)
)

func toolStatsQueryHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	accountID, projectID, err := GetToolReqTenant(tcr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query, err := GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	time, err := GetToolReqParam[string](tcr, "time", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.SelectAPIURL("stats_query"), nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	req.Header.Set("AccountID", accountID)
	req.Header.Set("ProjectID", projectID)

	q := req.URL.Query()
	q.Add("query", query)
	if time != "" {
		q.Add("time", time)
	}
	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolStatsQuery(s *server.MCPServer, c *config.Config) {
	s.AddTool(toolStatsQuery, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolStatsQueryHandler(ctx, c, request)
	})
}
