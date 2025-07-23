package tools_test

import (
	"context"
	"testing"

	"github.com/jgsheppa/everyday-mcp-server/pkg/tools"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestFrenchGreetingTool(t *testing.T) {
	ctx := context.Background()
	params := &mcp.CallToolParamsFor[tools.FrenchGreetingArgs]{
		Arguments: tools.FrenchGreetingArgs{Name: "Marie"},
	}

	frenchGreetingConfig := tools.NewFrenchGreetingTool()

	result, err := frenchGreetingConfig.ToolCall(ctx, nil, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Bonjour Marie!"
	if result.Content[0].(*mcp.TextContent).Text != expected {
		t.Errorf("expected %q, got %q", expected, result.Content[0].(*mcp.TextContent).Text)
	}
}

func TestFrenchGreetingToolEmptyName(t *testing.T) {
	ctx := context.Background()
	params := &mcp.CallToolParamsFor[tools.FrenchGreetingArgs]{
		Arguments: tools.FrenchGreetingArgs{Name: "  "},
	}

	frenchGreetingConfig := tools.NewFrenchGreetingTool()

	_, err := frenchGreetingConfig.ToolCall(ctx, nil, params)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

