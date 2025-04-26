package middleware

import "github.com/mark3labs/mcp-go/server"

type Middleware func(next server.ToolHandlerFunc) server.ToolHandlerFunc
