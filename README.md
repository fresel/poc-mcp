# poc-mcp

A small Model Context Protocol (MCP) playground with Go and TypeScript servers, plus a Go client.

## Repository layout

- `go-mcp-server/`: Go MCP server
- `go-mcp-client/`: Go stdio MCP client
- `ts-mcp-server/`: TypeScript MCP server
- `.vscode/mcp.json`: VS Code MCP server configuration
- `bin/`: locally built binaries

## Go MCP server

### Prerequisites

- Go 1.27 or later
- VS Code with MCP support, if using the server from VS Code

Build the server binary from the repository root:

```bash
mkdir -p bin
go build -o bin/go-mcp-server ./go-mcp-server
```

Run the server directly:

```bash
./bin/go-mcp-server
```

The server uses stdio, so it waits for MCP messages on standard input and writes responses to standard output.

### Configure VS Code

Create or update `.vscode/mcp.json`:

```json
{
  "servers": {
    "go-mcp-server": {
      "type": "stdio",
      "command": "${workspaceFolder}/bin/go-mcp-server"
    }
  }
}
```

After rebuilding the binary, restart or reload the MCP server from VS Code.

For the full setup and troubleshooting steps, see [docs/go-mcp-server.md](docs/go-mcp-server.md).

## Go checks

Run tests for the server:

```bash
cd go-mcp-server
go test ./...
```

Run tests for the client:

```bash
cd go-mcp-client
go test ./...
```

## Scaffold tool

The Go MCP server exposes `scaffold_java_app`. For example:

```json
{
  "folder": "demo",
  "kind": "controller"
}
```

This creates a Java scaffold under `test/demo/java-test`.
