set -e
set -o pipefail

go build -o ./bin/mcp-victorialogs ./cmd/mcp-victorialogs/main.go
