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
	toolHits = mcp.NewTool("hits",
		mcp.WithDescription("The number of matching log entries for the given <query> on the given [<start> ... <end>] time range grouped by <step> buckets. The returned results are sorted by time. This tool uses `/select/logsql/hits` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Log entries hits",
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
			mcp.Description(`LogsQL expression for the query of the logs data`),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Title("Start timestamp"),
			mcp.Description("Start timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("end",
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("step",
			mcp.Title("Step"),
			mcp.Description("The step is used to group the log entries by time. For example, 60s means that the log entries counts will be grouped by 1 minute bucket. Default is 1d."),
			mcp.Pattern(`^([0-9]+)([a-z]+)$`),
			mcp.DefaultString("1d"),
		),
		mcp.WithArray("field",
			mcp.Title("Field"),
			mcp.Description("Additionally, any number of field=<field_name> args can be passed to /select/logsql/hits for grouping hits buckets by the mentioned <field_name> fields. The grouped fields are put inside 'fields' object of response."),
		),
		mcp.WithNumber("fields_limit",
			mcp.Title("fields_limit"),
			mcp.Description("Optional fields_limit=N query arg can be passed for limiting the number of unique 'fields' groups to return to N. If more than N unique 'fields' groups is found, then top N 'fields' groups with the maximum number of 'total' hits are returned. The remaining hits are returned in 'fields': {} group."),
		),
	)
)

func toolHitsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	step, err := GetToolReqParam[string](tcr, "step", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	fields, err := GetToolReqParam[[]string](tcr, "field", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	fieldsLimit, err := GetToolReqParam[float64](tcr, "fields_limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.SelectAPIURL("hits"), nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	req.Header.Set("AccountID", accountID)
	req.Header.Set("ProjectID", projectID)

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", start)
	if end != "" {
		q.Add("end", end)
	}
	if step != "" {
		q.Add("step", step)
	}
	if len(fields) > 0 {
		for _, field := range fields {
			if field != "" {
				q.Add("field", field)
			}
		}
		if fieldsLimit > 0 {
			q.Add("fields_limit", fmt.Sprintf("%.f", fieldsLimit))
		}
	}
	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolHits(s *server.MCPServer, c *config.Config) {
	s.AddTool(toolHits, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolHitsHandler(ctx, c, request)
	})
}
