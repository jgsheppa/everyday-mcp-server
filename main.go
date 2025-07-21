package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type HiArgs struct {
	Name string `json:"name" jsonschema:"the name to say hi to"`
}

func greetTool(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[HiArgs]) (*mcp.CallToolResult, error) {
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

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "everyday",
		Version: "1.0.0",
	}, nil)

	greetToolDef := &mcp.Tool{
		Name:        "moin",
		Description: "Greets a person with the German \"Moin moin\" by name",
	}

	mcp.AddTool(server, greetToolDef, greetTool)

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}