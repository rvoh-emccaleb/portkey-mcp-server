# portkey-mcp-server

A Model Control Protocol (MCP) server implementation for Portkey. This application serves as a bridge for connecting various AI tools and services to Portkey through the Model Control Protocol.

## Table of Contents
- [Supported MCP Features](#supported-mcp-features)
    - [Tools](#tools)
- [Installation](#installation)
  - [Docker](#docker)
  - [Binary](#binary)
  - [From Source](#from-source)
- [Configuration](#configuration)
- [Usage](#usage)
  - [With Cursor IDE](#with-cursor-ide)
  - [With Claude Desktop](#with-claude-desktop)
  - [Manual SSE Requests](#manual-sse-requests)
- [Contributing](#contributing)
- [License](#license)

## Supported MCP Features
### Tools
- [`prompt_render`](https://portkey.ai/docs/api-reference/inference-api/prompts/render)

## Installation

### Docker

The easiest way to get started is with the pre-built Docker image:

```shell
# Run the container with the latest image (PORTKEY_API_KEY is required)
docker run -p 8080:8080 \
  -e TRANSPORT=sse \
  -e TRANSPORT_SSE_ADDRESS=:8080 \
  -e PORTKEY_API_KEY=your-api-key \
  ericmccaleb/portkey-mcp-server:latest
```

You can also build and run it locally using the Make targets:

```shell
# Clone the repository
git clone https://github.com/rvoh-emccaleb/portkey-mcp-server.git
cd portkey-mcp-server

# Build the Docker image locally (optional)
make docker-build

# Run the container (PORTKEY_API_KEY is required)
make docker-run PORTKEY_API_KEY=your-api-key

# Run with additional non-default settings
make docker-run PORT=9000 TRANSPORT=sse PORTKEY_API_KEY=your-api-key
```

You can also use the Docker commands directly:

```shell
# From repo root

# Build the Docker image
docker build --build-arg APP_VERSION=$(git rev-parse --short HEAD) -t portkey-mcp-server .

# Run in SSE mode (HTTP server)
docker run -p 8080:8080 \
  -e TRANSPORT=sse \
  -e TRANSPORT_SSE_ADDRESS=:8080 \
  -e PORTKEY_API_KEY=your-api-key \
  portkey-mcp-server
```

### Binary

To build a standalone binary:

```shell
# From repo root
make build
```

This will create a binary named `portkey-mcp-server` in the root directory that you can execute directly.

To run the binary, you'll need to provide the required environment variables:

```shell
# From repo root

# Run in SSE mode (HTTP server)
PORTKEY_API_KEY=your-api-key \
TRANSPORT=sse \
TRANSPORT_SSE_ADDRESS=:8080 \
./portkey-mcp-server

# Run in stdio mode
PORTKEY_API_KEY=your-api-key \
TRANSPORT=stdio \
./portkey-mcp-server
```

### From Source

To run the app directly with Go tooling:

```shell
# From repo root
go run -ldflags="-X main.appVersion=$(git rev-parse --short HEAD)" cmd/portkey-mcp-server/main.go
```

## Configuration

The server supports different transport configurations:
- `stdio` (Standard Input/Output) transport for command-line interfaces
- `sse` (Server-Sent Events) transport for HTTP-based communication

Configuration is handled through environment variables loaded at startup. [`.env.example`](./.env.example) can be used as a reference, but the source-of-truth is always the [config package](./internal/config/).

For running outside of Docker, you can configure the application by creating a `.env` file based on the variables expected by the [config package](./internal/config/). For Docker, environment variables should be set by other means.

## Usage

### With Cursor IDE

Create a `.cursor/mcp.json` file in your home directory (for global-level config) or in your repo root (for project-level config).

#### Using stdio Mode

With binary:
```json
{
  "mcpServers": {
    "Portkey": {
      "command": "/path/to/portkey-mcp-server-binary",
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "stdio"
      }
    }
  }
}
```

Or with Docker:
```json
{
  "mcpServers": {
    "Portkey": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-e", "PORTKEY_API_KEY",
        "-e", "TRANSPORT",
        "ericmccaleb/portkey-mcp-server:latest"
      ],
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "stdio"
      }
    }
  }
}
```

#### Using SSE Mode
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

Or to start up via Docker:
```json
{
  "mcpServers": {
    "Portkey": {
      "command": "docker",
      "args": [
        "run",
        "-p", "8080:8080",
        "--rm",
        "-e", "PORTKEY_API_KEY",
        "-e", "TRANSPORT",
        "-e", "TRANSPORT_SSE_ADDRESS",
        "ericmccaleb/portkey-mcp-server:latest"
      ],
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "sse",
        "TRANSPORT_SSE_ADDRESS": ":8080"
      }
    }
  }
}
```

### With Claude Desktop

Claude Desktop only supports stdio mode. Create a `claude_desktop_config.json` file in:
- macOS: `~/Library/Application Support/Claude/`
- Windows: `%APPDATA%\Claude\`

With binary:
```json
{
  "mcpServers": {
    "Portkey": {
      "command": "/path/to/portkey-mcp-server-binary",
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "stdio"
      }
    }
  }
}
```

Or with Docker:
```json
{
  "mcpServers": {
    "Portkey": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-e", "PORTKEY_API_KEY",
        "-e", "TRANSPORT",
        "ericmccaleb/portkey-mcp-server:latest"
      ],
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "stdio"
      }
    }
  }
}
```

### Manual SSE Requests

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

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on how to get started.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
