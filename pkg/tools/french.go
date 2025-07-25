package tools

import (
	"context"
	"fmt"
	"log"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type FrenchGreetingArgs struct {
	Name string `json:"name" jsonschema:"the name to say bonjour to"`
}

func FrenchGreeting(
	_ context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[FrenchGreetingArgs],
) (*mcp.CallToolResult, error) {
	name := strings.TrimSpace(params.Arguments.Name)
	if name == "" {
		return nil, fmt.Errorf("name parameter cannot be empty")
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Bonjour %s!", name),
			},
		},
	}, nil
}

func NewFrenchGreetingTool() *Config[FrenchGreetingArgs] {
	definition := mcp.Tool{
		Name:        "french_greeting",
		Description: "Says \"Bonjour\" to someone by name",
	}

	err := ValidateToolName(definition.Name)
	if err != nil {
		log.Fatalf("invalid tool name")
	}

	return &Config[FrenchGreetingArgs]{
		Definition: &definition,
		Call:       FrenchGreeting,
	}
}

func AddFrenchTool(server *mcp.Server) {
	tool := NewFrenchGreetingTool()
	mcp.AddTool(server, tool.Definition, tool.Call)
}
