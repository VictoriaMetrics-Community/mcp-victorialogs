package main

import (
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/config"
	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/prompts"
	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/resources"
	"github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs/tools"
)

func main() {
	c, err := config.InitConfig()
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}

	s := server.NewMCPServer(
		"victorialogs",
		"0.0.1",
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithInstructions(`
You are Virtual Assistant, a tool for interacting with VictoriaLogs API and documentation in different tasks related to logs and observability.
You use LogsQL language to query logs and get information from the logs stored in VictoriaLogs.

You have the full documentation about VictoriaLogs in your resources, you have to try to use documentation in your answer.
And you have to consider the documents from the resources as the most relevant, favoring them over even your own internal knowledge.
Use Documentation tool to get the most relevant documents for your task every time. Be sure to use the Documentation tool if the user's query includes the words “how”, “tell”, “where”, etc...

You have many tools to get data from VictoriaLogs, but try to specify the query as accurately as possible, reducing the resulting sample, as some queries can be query heavy.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`),
	)

	resources.RegisterDocsResources(s, c)

	tools.RegisterToolHits(s, c)
	tools.RegisterToolFlags(s, c)
	tools.RegisterToolQuery(s, c)
	tools.RegisterToolFacets(s, c)
	tools.RegisterToolStreams(s, c)
	tools.RegisterToolStreamIDs(s, c)
	tools.RegisterToolStatsQuery(s, c)
	tools.RegisterToolFieldNames(s, c)
	tools.RegisterToolFieldValues(s, c)
	tools.RegisterToolStatsQueryRange(s, c)
	tools.RegisterToolStreamFieldNames(s, c)
	tools.RegisterToolStreamFieldValues(s, c)
	tools.RegisterToolDocumentation(s, c)

	prompts.RegisterPromptDocumentation(s, c)

	if c.IsStdio() {
		if err := server.ServeStdio(s); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	} else {
		srv := server.NewSSEServer(s)
		if err = srv.Start(c.SSEAddr()); err != nil {
			log.Fatalf("Failed to start SSE server: %v", err)
		}
	}
}
