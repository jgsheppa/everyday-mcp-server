package tools

import (
	"context"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Config[T any] struct {
	Definition *mcp.Tool
	Call       func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[T]) (*mcp.CallToolResult, error)
}
