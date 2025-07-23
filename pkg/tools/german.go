package tools

import (
	"context"
	"fmt"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type GermanGreetingArgs struct {
	Name string `json:"name" jsonschema:"the name to say moin moin to"`
}

func GermanGreetingTool(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[GermanGreetingArgs]) (*mcp.CallToolResult, error) {
	name := strings.TrimSpace(params.Arguments.Name)
	if name == "" {
		return nil, fmt.Errorf("name parameter cannot be empty")
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Moin moin %s!", name),
			},
		},
	}, nil
}

func NewGermanGreetingTool() *Config[GermanGreetingArgs] {
	return &Config[GermanGreetingArgs]{Definition: &mcp.Tool{
		Name:        "German greeting",
		Description: "Says \"Moin moin\" to someone by name"},
		ToolCall: GermanGreetingTool,
	}
}

func AddGermanTool(server *mcp.Server) {
	tool := NewGermanGreetingTool()
	mcp.AddTool(server, tool.Definition, tool.ToolCall)
}
