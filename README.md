# VictoriaLogs MCP Server

[![smithery badge](https://smithery.ai/badge/@VictoriaMetrics-Community/mcp-victorialogs)](https://smithery.ai/server/@VictoriaMetrics-Community/mcp-victorialogs)

The implementation of [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/).

This provides access to your VictoriaLogs instance and seamless integration with [VictoriaLogs APIs](https://docs.victoriametrics.com/victorialogs/querying/#http-api) and [documentation](https://docs.victoriametrics.com/victorialogs/).
It can give you a comprehensive interface for logs, observability, and debugging tasks related to your VictoriaLogs instances, enable advanced automation and interaction capabilities for engineers and tools.

## Features

This MCP server allows you to use almost all read-only APIs of VictoriaLogs, i.e. all functions available in [Web UI](https://docs.victoriametrics.com/victorialogs/querying/#web-ui):

- Querying logs and exploring logs data
- Showing parameters of your VictoriaLogs instance
- Listing available streams, fields, field values
- Query statistics for the logs as metrics
 
In addition, the MCP server contains embedded up-to-date documentation and is able to search it without online access.

More details about the exact available tools and prompts can be found in the [Usage](#usage) section.

You can combine functionality of tools, docs search in your prompts and invent great usage scenarios for your VictoriaLogs instance.
And please note the fact that the quality of the MCP Server and its responses depends very much on the capabilities of your client and the quality of the model you are using.

You can also combine the MCP server with other observability or doc search related MCP Servers and get even more powerful results.

## Requirements

- [VictoriaLogs](https://docs.victoriametrics.com/victorialogs/) instance ([single-node](https://docs.victoriametrics.com/victorialogs/) or [cluster](https://docs.victoriametrics.com/victorialogs/cluster/))
- Go 1.24 or higher (if you want to build from source)

## Installation

### Go

```bash
go install github.com/VictoriaMetrics-Community/mcp-victorialogs/cmd/mcp-victorialogs@latest
```

### Source Code

```bash
git clone https://github.com/VictoriaMetrics-Community/mcp-victorialogs.git
cd mcp-victorialogs
go build -o bin/mcp-victorialogs ./cmd/mcp-victorialogs/main.go

# after that add bin/mcp-victorialogs file to your PATH
```

### Binaries

Just download the latest release from [Releases](https://github.com/VictoriaMetrics-Community/mcp-victorialogs/releases) page and put it to your PATH.

### Docker

Coming soon...

### Smithery

To install VictoriaLogs MCP Server for your client automatically via [Smithery](https://smithery.ai/server/@VictoriaMetrics-Community/mcp-victorialogs), yo can use the following commands:

```bash
# Get the list of supported MCP clients
npx -y @smithery/cli list clients
#Available clients:
#  claude
#  cline
#  windsurf
#  roocode
#  witsy
#  enconvo
#  cursor
#  vscode
#  vscode-insiders
#  boltai
#  amazon-bedrock

# Install VictoriaLogs MCP server for your client
npx -y @smithery/cli install @VictoriaMetrics-Community/mcp-victorialogs --client <YOUR-CLIENT-NAME>
# and follow the instructions
```

## Configuration

MCP Server for VictoriaLogs is configured via environment variables:

| Variable | Description                               | Required | Default | Allowed values |
|----------|-------------------------------------------|----------|---------|---------|
| `VL_INSTANCE_ENTRYPOINT` | URL to VictoriaLogs instance              | Yes | -       | - |
| `VL_INSTANCE_BEARER_TOKEN` | Authentication token for VictoriaLogs API | No | -       | - |
| `MCP_SERVER_MODE` | Server operation mode                     | No | `stdio` | `stdio`, `sse` |
| `MCP_SSE_ADDR` | Address for SSE server to listen on       | No | `:8081` | - |

### Сonfiguration examples

```bash
# For a public playground
export VL_INSTANCE_ENTRYPOINT="https://play-vmlogs.victoriametrics.com"

# Server mode
export MCP_SERVER_MODE="sse"
export MCP_SSE_ADDR="0.0.0.0:8081"
```

## Setup in clients

### Cursor

Go to: `Settings` -> `Cursor Settings` -> `MCP` -> `Add new global MCP server` and paste the following configuration into your Cursor `~/.cursor/mcp.json` file:

```json
{
  "mcpServers": {
    "victorialogs": {
      "command": "/path/to/mcp-victorialogs",
      "env": {
        "VL_INSTANCE_ENTRYPOINT": "<YOUR_VL_INSTANCE>",
        "VL_INSTANCE_BEARER_TOKEN": "<YOUR_VL_BEARER_TOKEN>"
      }
    }
  }
}
```

See [Cursor MCP docs](https://docs.cursor.com/context/model-context-protocol) for more info.

### Claude Desktop

Add this to your Claude Desktop `claude_desktop_config.json` file (you can find it if open `Settings` -> `Developer` -> `Edit config`):

```json
{
  "mcpServers": {
    "victorialogs": {
      "command": "/path/to/mcp-victorialogs",
      "env": {
        "VL_INSTANCE_ENTRYPOINT": "<YOUR_VL_INSTANCE>",
        "VL_INSTANCE_BEARER_TOKEN": "<YOUR_VL_BEARER_TOKEN>"
      }
    }
  }
}
```

See [Claude Desktop MCP docs](https://modelcontextprotocol.io/quickstart/user) for more info.

### Claude Code

Run the command:

```sh
claude mcp add victorialogs -- /path/to/mcp-victorialogs \
  -e VL_INSTANCE_ENTRYPOINT=<YOUR_VL_INSTANCE> \
  -e VL_INSTANCE_BEARER_TOKEN=<YOUR_VL_BEARER_TOKEN>
```

See [Claude Code MCP docs](https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/tutorials#set-up-model-context-protocol-mcp) for more info.

### Visual Studio Code

Add this to your VS Code MCP config file:

```json
{
  "servers": {
    "victorialogs": {
      "type": "stdio",
      "command": "/path/to/mcp-victorialogs",
      "env": {
        "VL_INSTANCE_ENTRYPOINT": "<YOUR_VL_INSTANCE>",
        "VL_INSTANCE_BEARER_TOKEN": "<YOUR_VL_BEARER_TOKEN>"
      }
    }
  }
}
```

See [VS Code MCP docs](https://code.visualstudio.com/docs/copilot/chat/mcp-servers) for more info.

### Zed

Add the following to your Zed config file:

```json
  "context_servers": {
    "victorialogs": {
      "command": {
        "path": "/path/to/mcp-victorialogs",
        "args": [],
        "env": {
          "VL_INSTANCE_ENTRYPOINT": "<YOUR_VL_INSTANCE>",
          "VL_INSTANCE_BEARER_TOKEN": "<YOUR_VL_BEARER_TOKEN>"
        }
      },
      "settings": {}
    }
  }
}
```

See [Zed MCP docs](https://zed.dev/docs/ai/mcp) for more info.

### JetBrains IDEs

- Open `Settings` -> `Tools` -> `AI Assistant` -> `Model Context Protocol (MCP)`.
- Click `Add (+)`
- Select `As JSON`
- Put the following to the input field:

```json
{
  "mcpServers": {
    "victorialogs": {
      "command": "/path/to/mcp-victorialogs",
      "env": {
        "VL_INSTANCE_ENTRYPOINT": "<YOUR_VL_INSTANCE>",
        "VL_INSTANCE_BEARER_TOKEN": "<YOUR_VL_BEARER_TOKEN>"
      }
    }
  }
}
```

### Windsurf

Add the following to your Windsurf MCP config file.

```json
{
  "mcpServers": {
    "victorialogs": {
      "command": "/path/to/mcp-victorialogs",
      "env": {
        "VL_INSTANCE_ENTRYPOINT": "<YOUR_VL_INSTANCE>",
        "VL_INSTANCE_BEARER_TOKEN": "<YOUR_VL_BEARER_TOKEN>"
      }
    }
  }
}
```

See [Windsurf MCP docs](https://docs.windsurf.com/windsurf/mcp) for more info.

### Amazon Bedrock

Coming soon....

### Using Docker instead of binary

Coming soon...

## Usage

After [installing](#installation) and [configuring](#setup-in-clients) the MCP server, you can start using it with your favorite MCP client.

You can start dialog with AI assistant from the phrase:

```
Use MCP VictoriaLogs in the following answers
```

But it's not required, you can just start asking questions and the assistant will automatically use the tools and documentation to provide you with the best answers.

### Toolset

MCP VictoriaLogs provides numerous tools for interacting with your VictoriaLogs instance.

Here's a list of available tools:

| Tool                 | Description                                           |
|----------------------|-------------------------------------------------------|
| `documentation`      | Search in embedded VictoriaLogs documentation         |
| `facets`             | Most frequent values per each log field               |
| `field_names`        | List of field names for the query                     |
| `field_values`       | List of field values for the query                    |
| `flags`              | View non-default flags of the VictoriaLogs instance   |
| `hits`               | The number of matching log entries grouped by buckets |
| `query`              | Execute LogsQL queries                                |
| `stats_query`        | Querying log stats for the given time                 |
| `stats_query_range`  | Querying log stats on the given time range            |
| `stream_field_names` | List of stream fields for the query                   |
| `stream_field_names` | List of stream field values for the query             |
| `stream_ids`         | List of stream IDs for the query                      |
| `streams`            | List of streams for the query                         |

### Prompts

The server includes pre-defined prompts for common tasks.

These are just examples at the moment, the prompt library will be added to in the future:

| Prompt | Description                                           |
|--------|-------------------------------------------------------|
| `documentation` | Search VictoriaLogs documentation for specific topics |

## Disclaimer

AI services and agents along with MCP servers like this cannot guarantee the accuracy, completeness and reliability of results.
You should double check the results obtained with AI.
The quality of the MCP Server and its responses depends very much on the capabilities of your client and the quality of the model you are using.

## Contributing

Contributions to the MCP VictoriaLogs project are welcome! Please feel free to submit issues, feature requests, or pull requests.
