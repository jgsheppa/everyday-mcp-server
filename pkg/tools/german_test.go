package tools_test

import (
	"context"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jgsheppa/everyday-mcp-server/pkg/tools"
)

func TestGermanGreetingTool(t *testing.T) {
	ctx := context.Background()
	params := &mcp.CallToolParamsFor[tools.GermanGreetingArgs]{
		Arguments: tools.GermanGreetingArgs{Name: "Alice"},
	}

	germanGreetingConfig := tools.NewGermanGreetingTool()

	result, err := germanGreetingConfig.Call(ctx, nil, params)
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

	_, err := germanGreetingConfig.Call(ctx, nil, params)
	if err == nil {
		t.Error("expected error for empty name")
	}
}
