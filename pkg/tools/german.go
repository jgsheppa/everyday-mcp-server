package tools

import (
	"context"
	"fmt"
	"log"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type GermanGreetingArgs struct {
	Name string `json:"name" jsonschema:"the name to say moin moin to"`
}

func GermanGreeting(
	_ context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GermanGreetingArgs],
) (*mcp.CallToolResult, error) {
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
	definition := mcp.Tool{
		Name:        "german_greeting",
		Description: "Says \"Moin moin\" to someone by name",
	}
	err := ValidateToolName(definition.Name)
	if err != nil {
		log.Fatalf("invalid tool name")
	}

	return &Config[GermanGreetingArgs]{
		Definition: &definition,
		Call:       GermanGreeting,
	}
}

func AddGermanTool(server *mcp.Server) {
	tool := NewGermanGreetingTool()
	mcp.AddTool(server, tool.Definition, tool.Call)
}
