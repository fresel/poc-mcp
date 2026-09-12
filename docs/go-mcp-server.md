# Configure and build the Go MCP server

This guide configures the Go MCP server as a VS Code stdio server and builds the Docker image that VS Code launches.

## 1. Check Go

From the repository root, confirm that Go is installed:

```bash
go version
```

The Go modules in this repository require Go 1.27 or later.

## 2. Build the binary

Run this command from the repository root:

```bash
docker build -f docker-mcp-server/Dockerfile -t poc-mcp-server:distroless .
```

The image is tagged as:

```text
poc-mcp-server:distroless
```


## 3. Configure VS Code

Add the following server entry to `.vscode/mcp.json`:

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
        "poc-mcp-server:distroless"
      ]
    }
  }
}
```
`${workspaceFolder}` resolves to the repository root opened in VS Code. The Docker volume makes that workspace available inside the container at `/workspace`.
`${workspaceFolder}` resolves to the repository root opened in VS Code. The command must point to the compiled binary, not to `go run`, because this configuration is intended to use the build in `bin/`.

If the workspace also contains other MCP servers, keep their entries under the same `servers` object.

## 4. Start or reload the server

In VS Code:

1. Open the Command Palette.
2. Run `MCP: List Servers`.
3. Select `go-mcp-server` and start or restart it.

The server communicates over standard input and standard output. Do not add logging to standard output, since it would corrupt the MCP protocol stream.

## Choose the scaffold destination

The `scaffold_java_app` tool accepts a `folder` argument supplied by the MCP user:

```json
{
  "folder": "/workspaces/my-project",
  "kind": "controller"
}
```

The folder can be an absolute path or a path relative to the server's working directory. The tool creates the `java-test` project inside it, so the example creates:

```text
/workspaces/my-project/java-test
```

Supported `kind` values are `controller`, `service`, `integration-test`, and `unit-test`. The target must not already exist.

When running the server in Docker, the selected host directory must be mounted into the container. For example, to allow scaffolding under the repository's `test` directory:

```json
{
  "go-mcp-server": {
    "type": "stdio",
    "command": "docker",
    "args": [
      "run",
      "--rm",
      "-i",
      "-v",
      "${workspaceFolder}/test:/workspace/test",
      "poc-mcp-server:distroless"
    ]
  }
}
```

Use `folder: "/workspace/test/demo"` with that configuration. It creates `/workspace/test/demo/java-test` in the container, persisted on the host as `test/demo/java-test`.

## 5. Rebuild after changes

Whenever Go server code changes, rebuild the binary:

```bash
docker build -f docker-mcp-server/Dockerfile -t poc-mcp-server:distroless .
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

Run the Docker build from the repository root:

```bash
docker build -f docker-mcp-server/Dockerfile -t poc-mcp-server:distroless .
```

Or change into `go-mcp-server` before running module commands.

### VS Code cannot start the server

Check that the image exists:

```bash
docker image inspect poc-mcp-server:distroless
```

Rebuild the image if necessary, then restart the server from VS Code.
