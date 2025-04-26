# portkey-mcp-server

A Model Control Protocol (MCP) server implementation for Portkey. This application serves as a bridge for connecting various AI tools and services to Portkey through the Model Control Protocol.

# Supported MCP Features
## Tools
- [`prompt_render`](https://portkey.ai/docs/api-reference/inference-api/prompts/render)

# Getting Started

### Local Development

Execute the following to install git hooks in your local repo, which will ensure that mocks are regenerated and committed before pushing:
```shell
cd <root-of-repo>
make install-hooks
```

If you are seeing stale linter errors coming from the result of `make lint` (part of those installed git hooks), you could try clearing your linter cache with `make lint-clear-cache`.

### Configuration

The server supports different transport configurations:
- `stdio` (Standard Input/Output) transport for command-line interfaces
- `sse` (Server-Sent Events) transport for HTTP-based communication

Configuration is handled through environment variables loaded at startup. [`.env.example`](./.env.example) can be used as a reference, but the source-of-truth is always the [config package](./internal/config/).

For running outside of Docker, you can configure the application by creating a `.env` file based on the variables expected by the [config package](./internal/config/). For Docker, environment variables should be set by other means.

### Running the app

#### Direct Go Execution

To run the app locally with the correct version information, you can execute:

```shell
cd <root-of-repo>
go run -ldflags="-X main.appVersion=$(git rev-parse --short HEAD)" cmd/portkey-mcp-server/main.go
```

#### Binary Execution

Alternatively, you can build a binary via:

```shell
make build
```

This will create a binary named `portkey-mcp-server` in the root directory that you can execute directly.

#### Running with Docker

The application can be built and run using Docker with the provided Make targets:

```shell
# Build the Docker image
make docker-build

# Run the container (PORTKEY_API_KEY is required)
make docker-run PORTKEY_API_KEY=your-api-key

# Run with additional non-default settings
make docker-run PORT=9000 TRANSPORT=sse PORTKEY_API_KEY=your-api-key
```

You can also use the Docker commands directly:

```shell
# Build the Docker image
docker build --build-arg APP_VERSION=$(git rev-parse --short HEAD) -t portkey-mcp-server .

# Run in SSE mode (HTTP server)
docker run -p 8080:8080 \
  -e TRANSPORT=sse \
  -e TRANSPORT_SSE_ADDRESS=:8080 \
  -e PORTKEY_API_KEY=your-api-key \
  portkey-mcp-server
```

### Interacting with the MCP Server

#### Accessing the SSE Server Manually

When running in SSE mode, you can access the server using HTTP requests. Following the MCP protocol requires these steps:

```shell
# 1. First establish the SSE connection to get a session
curl http://localhost:8080/sse
# This will return something like:
# event: endpoint
# data: /message?sessionId=abc123

# 2. Initialize the session (using the URL from #1, above)
curl -X POST "http://localhost:8080/message?sessionId=abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "1",
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "clientInfo": {
        "name": "curl-client",
        "version": "1.0.0"
      },
      "capabilities": {}
    }
  }'

# 3. Now you can make tool calls
curl -X POST "http://localhost:8080/message?sessionId=abc123" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "2",
    "method": "call_tool",
    "params": {
      "name": "prompt_render", 
      "arguments": {
        "prompt_id": "your_prompt_id",
        "variables": {
            "your_prompt_variable": "some_value"
        }
      }
    }
  }'
```

> Note: The SSE connection in step 1 must remain open during your session. Run it in a separate terminal window and do not close it until you're done with the session.

#### Using with Cursor IDE

To use this MCP server with Cursor IDE, you'll need to configure a `.cursor/mcp.json` file:

##### Using SSE Mode

If the server is already running locally:
```json
{
  "mcpServers": {
    "Portkey": {
      "url": "http://localhost:8080/sse"
    }
  }
}
```

To start up via Docker:
```json
{
  "mcpServers": {
    "Portkey": {
      "command": "docker",
      "args": [
        "run",
        "-p",
        "8080:8080",
        "--rm",
        "-e",
        "TRANSPORT",
        "-e",
        "TRANSPORT_SSE_ADDRESS",
        "-e",
        "PORTKEY_API_KEY",
        "portkey-mcp-server:latest"
      ],
      "env": {
        "TRANSPORT": "sse",
        "TRANSPORT_SSE_ADDRESS": ":8080",
        "PORTKEY_API_KEY": "your-api-key"
      }
    }
  }
}
```

##### Using stdio Mode (Local Binary)

```json
{
  "mcpServers": {
    "Portkey": {
      "command": "/path/to/portkey-mcp-server-binary",
      "env": {
        "TRANSPORT": "stdio",
        "PORTKEY_API_KEY": "your-api-key"
      }
    }
  }
}
```

### Running Tests

This project uses various testing tools:

```shell
# Run all tests
make test

# Run benchmarks
make benchmark

# Generate mocks
make mocks

# Run linter
make lint

# Run security checks
make security
```
