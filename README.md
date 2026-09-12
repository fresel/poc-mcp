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

Build the Docker image from the repository root:

```bash
docker build -f docker-mcp-server/Dockerfile -t poc-mcp-server:distroless .
```

Run the server directly:

```bash
docker run --rm -i poc-mcp-server:distroless
```

The server uses stdio, so it waits for MCP messages on standard input and writes responses to standard output.

### Configure VS Code

Create or update `.vscode/mcp.json`:

```json
{
  "servers": {
    "go-mcp-server": {
      "type": "stdio",
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-v",
        "${workspaceFolder}:/workspace",
        "-e",
        "MCP_WORKSPACE_ROOT=/workspace",
        "poc-mcp-server:distroless"
      ]
    }
  }
}
```
After rebuilding the image, restart or reload the MCP server from VS Code.

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
  "projectName": "my-project",
  "folder": "demo",
  "kind": "controller"
}
```

`folder` is relative to the workspace root and `projectName` is the project directory name. The example creates `demo/my-project`. Call the tool again with the same `projectName` and another `kind` (`controller`, `service`, `model`, `integration-test`, or `unit-test`) to add another scaffold to that project.
