package main

import (
	"context"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Add to main_test.go
func TestGreetTool(t *testing.T) {
    ctx := context.Background()
    params := &mcp.CallToolParamsFor[HiArgs]{
        Arguments: HiArgs{Name: "Alice"},
    }
    
    result, err := greetTool(ctx, nil, params)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    expected := "Moin moin Alice!"
    if result.Content[0].(*mcp.TextContent).Text != expected {
        t.Errorf("expected %q, got %q", expected, result.Content[0].(*mcp.TextContent).Text)
    }
}

func TestGreetToolEmptyName(t *testing.T) {
    ctx := context.Background()
    params := &mcp.CallToolParamsFor[HiArgs]{
        Arguments: HiArgs{Name: "  "},
    }
    
    _, err := greetTool(ctx, nil, params)
    if err == nil {
        t.Error("expected error for empty name")
    }
}