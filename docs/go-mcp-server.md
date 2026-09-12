# Configure and build the Go MCP server

This guide configures the Go MCP server as a VS Code stdio server and builds the binary that VS Code launches.

## 1. Check Go

From the repository root, confirm that Go is installed:

```bash
go version
```

The Go modules in this repository require Go 1.27 or later.

## 2. Build the binary

Run this command from the repository root:

```bash
mkdir -p bin
go build -o bin/go-mcp-server ./go-mcp-server
```

The output is:

```text
bin/go-mcp-server
```

The `bin/` directory is ignored by Git because the binary is a local build artifact.

To verify the binary exists and is executable:

```bash
test -x bin/go-mcp-server && file bin/go-mcp-server
```

## 3. Configure VS Code

Add the following server entry to `.vscode/mcp.json`:

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

`${workspaceFolder}` resolves to the repository root opened in VS Code. The command must point to the compiled binary, not to `go run`, because this configuration is intended to use the build in `bin/`.

If the workspace also contains other MCP servers, keep their entries under the same `servers` object.

## 4. Start or reload the server

In VS Code:

1. Open the Command Palette.
2. Run `MCP: List Servers`.
3. Select `go-mcp-server` and start or restart it.

The server communicates over standard input and standard output. Do not add logging to standard output, since it would corrupt the MCP protocol stream.

## 5. Rebuild after changes

Whenever Go server code changes, rebuild the binary:

```bash
go build -o bin/go-mcp-server ./go-mcp-server
```

Then restart the MCP server in VS Code.

## 6. Validate the server module

Run the Go tests from the repository root:

```bash
go test ./go-mcp-server/...
```

You can also build from inside the server module, writing to the repository-level `bin/` directory:

```bash
cd go-mcp-server
mkdir -p ../bin
go build -o ../bin/go-mcp-server .
```

## Troubleshooting

### `go: cannot find main module`

Run the command from the repository root with the package path:

```bash
go build -o bin/go-mcp-server ./go-mcp-server
```

Or change into `go-mcp-server` before running module commands.

### VS Code cannot start the server

Check that the binary exists and is executable:

```bash
ls -l bin/go-mcp-server
test -x bin/go-mcp-server
```

Rebuild it if necessary, then restart the server from VS Code.
