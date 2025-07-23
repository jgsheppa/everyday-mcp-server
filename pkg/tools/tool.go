package tools

import (
	"context"
	"fmt"
	"regexp"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var toolNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{1,64}$`)

type Config[T any] struct {
	Definition *mcp.Tool
	Call       func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[T]) (*mcp.CallToolResult, error)
}

func ValidateToolName(name string) error {
	if !toolNameRegex.MatchString(name) {
		return fmt.Errorf("tool name '%s' must match regex [a-zA-Z0-9_]{1,64}", name)
	}
	return nil
}
