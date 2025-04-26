package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Tuple represents a tool and its handler function.
// This is used to register tools with the MCP server.
type Tuple struct {
	Tool    *mcp.Tool
	Handler server.ToolHandlerFunc
}
