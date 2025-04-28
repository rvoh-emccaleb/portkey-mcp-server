# portkey-mcp-server

A Model Control Protocol (MCP) server implementation for Portkey. This application serves as a bridge for connecting various AI tools and services to Portkey through the Model Control Protocol.

## Table of Contents
- [Supported MCP Features](#supported-mcp-features)
- [Getting Started](#getting-started)
  - [Configuration](#configuration)
  - [Running the app](#running-the-app)
    - [Running with Docker](#running-with-docker)
    - [Binary Execution](#binary-execution)
    - [Direct Go Execution](#direct-go-execution)
- [Interacting with the MCP Server](#interacting-with-the-mcp-server)
  - [Using with Cursor IDE or Claude Desktop](#using-with-cursor-ide-or-claude-desktop)
  - [Accessing the SSE Server Manually](#accessing-the-sse-server-manually)
- [Contributing](#contributing)
  - [Local Development](#local-development)
  - [Submitting a Pull Request](#submitting-a-pull-request)
    - [Running Tests](#running-tests)
    - [CI/CD Requirements](#cicd-requirements)
  - [License](#license)

# Supported MCP Features
## Tools
- [`prompt_render`](https://portkey.ai/docs/api-reference/inference-api/prompts/render)

# Getting Started

### Configuration

The server supports different transport configurations:
- `stdio` (Standard Input/Output) transport for command-line interfaces
- `sse` (Server-Sent Events) transport for HTTP-based communication

Configuration is handled through environment variables loaded at startup. [`.env.example`](./.env.example) can be used as a reference, but the source-of-truth is always the [config package](./internal/config/).

For running outside of Docker, you can configure the application by creating a `.env` file based on the variables expected by the [config package](./internal/config/). For Docker, environment variables should be set by other means.

### Running the app

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

#### Binary Execution

Alternatively, you can build a binary via:

```shell
make build
```

This will create a binary named `portkey-mcp-server` in the root directory that you can execute directly.

#### Direct Go Execution

To run the app with Go tooling, you can execute:

```shell
cd <root-of-repo>
go run -ldflags="-X main.appVersion=$(git rev-parse --short HEAD)" cmd/portkey-mcp-server/main.go
```

### Interacting with the MCP Server

#### Using with Cursor IDE or Claude Desktop

The MCP server can be configured for different clients and transport modes. Choose the configuration that matches your use case:

##### For Claude Desktop (stdio mode only)

Create a `claude_desktop_config.json` file in:
- macOS: `~/Library/Application Support/Claude/`
- Windows: `%APPDATA%\Claude\`

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

Or using Docker:
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
        "portkey-mcp-server:latest"
      ],
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "stdio"
      }
    }
  }
}
```

##### For Cursor IDE

Create a `.cursor/mcp.json` file in your home directory (for global-level config) or in your repo root (for project-level config).

###### Cursor: Using stdio Mode
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

Or using Docker:
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
        "portkey-mcp-server:latest"
      ],
      "env": {
        "PORTKEY_API_KEY": "your-api-key",
        "TRANSPORT": "stdio"
      }
    }
  }
}
```

###### Cursor: Using SSE Mode
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
        "portkey-mcp-server:latest"
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

## Contributing

We welcome contributions to this project! Here's how to get started:

### Local Development

Execute the following to install git hooks in your local repo, which will ensure that mocks are regenerated and committed before pushing:
```shell
cd <root-of-repo>
make install-hooks
```

If you are seeing stale linter errors coming from the result of `make lint` (part of those installed git hooks), you could try clearing your linter cache with `make lint-clear-cache`.

### Submitting a Pull Request

1. Create a new branch for your changes
2. Make your changes and commit them with clear, descriptive commit messages
3. Run all tests locally using the commands above
4. Push your branch and create a Pull Request

#### Running Tests

Before submitting a PR, please consider running all tests locally to ensure your changes don't introduce any issues:

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

> Note: The above are run automatically in our GitHub Actions workflows. Running them locally before pushing reduces noise in your PR and helps catch issues early.

#### CI/CD Requirements

All GitHub Actions jobs must pass before a PR can be merged. This includes:
- Build verification
- Unit tests
- Linting
- Security checks
- etc.

If any job fails, you can find detailed error output in the GitHub Actions artifacts:
1. Go to the "Actions" tab in the repository
2. Click on the failed workflow run
3. Click on the failed job
4. Look for the "Artifacts" section at the bottom of the job page
5. Download and review the relevant artifacts (e.g., `lint-report.json` for linter errors)

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
