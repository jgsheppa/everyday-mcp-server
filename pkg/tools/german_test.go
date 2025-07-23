package tools_test

import (
	"context"
	"testing"

	"github.com/jgsheppa/everyday-mcp-server/pkg/tools"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGermanGreetingTool(t *testing.T) {
	ctx := context.Background()
	params := &mcp.CallToolParamsFor[tools.GermanGreetingArgs]{
		Arguments: tools.GermanGreetingArgs{Name: "Alice"},
	}

	germanGreetingConfig := tools.NewGermanGreetingTool()

	result, err := germanGreetingConfig.ToolCall(ctx, nil, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Moin moin Alice!"
	if result.Content[0].(*mcp.TextContent).Text != expected {
		t.Errorf("expected %q, got %q", expected, result.Content[0].(*mcp.TextContent).Text)
	}
}

func TestGermanGreetingToolEmptyName(t *testing.T) {
	ctx := context.Background()
	params := &mcp.CallToolParamsFor[tools.GermanGreetingArgs]{
		Arguments: tools.GermanGreetingArgs{Name: "  "},
	}

	germanGreetingConfig := tools.NewGermanGreetingTool()

	_, err := germanGreetingConfig.ToolCall(ctx, nil, params)
	if err == nil {
		t.Error("expected error for empty name")
	}
}
