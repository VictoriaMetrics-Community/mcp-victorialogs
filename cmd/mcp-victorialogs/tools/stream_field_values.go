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
	toolStreamFieldValues = mcp.NewTool("stream_field_values",
		mcp.WithDescription("Get unique values for the given <fieldName> field from results of the given <query> on the given [<start> ... <end>] time range. The response also contains the number of log results per every field value. This tool uses `/select/logsql/stream_field_values` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of stream field values",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
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
			mcp.Description(`LogsQL expression for the query of the logs streams`),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Title("Start timestamp"),
			mcp.Description("Start timestamp in RFC3339"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("end",
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in RFC3339. If <end> is missing, then it equals to the maximum timestamp across logs stored in VictoriaLogs."),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("field",
			mcp.Required(),
			mcp.Title("Field name"),
			mcp.Description("Field name of log stream for getting field values."),
		),
	)
)

func toolStreamFieldValuesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	accountID, projectID, err := GetToolReqTenant(tcr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query, err := GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start, err := GetToolReqParam[string](tcr, "start", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	end, err := GetToolReqParam[string](tcr, "end", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	field, err := GetToolReqParam[string](tcr, "field", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.SelectAPIURL("stream_field_values"), nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	req.Header.Set("AccountID", accountID)
	req.Header.Set("ProjectID", projectID)

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", start)
	q.Add("field", field)
	if end != "" {
		q.Add("end", end)
	}
	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolStreamFieldValues(s *server.MCPServer, c *config.Config) {
	s.AddTool(toolStreamFieldValues, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolStreamFieldValuesHandler(ctx, c, request)
	})
}
