package tools

import (
	"context"
	"testing"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestYearsSince(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		expectError bool
	}{
		{
			name:        "valid past year",
			year:        2020,
			expectError: false,
		},
		{
			name:        "current year",
			year:        2025,
			expectError: false,
		},
		{
			name:        "future year",
			year:        2030,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &mcp.CallToolParamsFor[YearsSinceArgs]{
				Arguments: YearsSinceArgs{
					Year: tt.year,
				},
			}

			result, err := YearsSince(context.Background(), nil, params)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("expected result but got nil")
				return
			}

			if len(result.Content) == 0 {
				t.Errorf("expected content but got none")
			}
		})
	}
}

func TestNewYearsSinceTool(t *testing.T) {
	tool := NewYearsSinceTool()
	
	if tool.Error != nil {
		t.Errorf("unexpected error: %v", tool.Error)
	}

	if tool.Definition == nil {
		t.Errorf("expected tool definition but got nil")
	}

	if tool.Definition.Name != "years_since" {
		t.Errorf("expected tool name 'years_since', got '%s'", tool.Definition.Name)
	}
}