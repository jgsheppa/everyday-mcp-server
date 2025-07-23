package tools_test

import (
	"strings"
	"testing"

	"github.com/jgsheppa/everyday-mcp-server/pkg/tools"
)

func TestValidateToolName(t *testing.T) {
	tests := []struct {
		name    string
		toolName string
		wantErr bool
	}{
		{
			name:     "valid tool name with underscores",
			toolName: "french_greeting",
			wantErr:  false,
		},
		{
			name:     "valid tool name with numbers",
			toolName: "tool123",
			wantErr:  false,
		},
		{
			name:     "valid tool name mixed case",
			toolName: "MyTool_123",
			wantErr:  false,
		},
		{
			name:     "valid single character",
			toolName: "a",
			wantErr:  false,
		},
		{
			name:     "valid 64 character name",
			toolName: strings.Repeat("a", 64),
			wantErr:  false,
		},
		{
			name:     "invalid empty name",
			toolName: "",
			wantErr:  true,
		},
		{
			name:     "invalid name with spaces",
			toolName: "French Greeting",
			wantErr:  true,
		},
		{
			name:     "invalid name with hyphens",
			toolName: "french-greeting",
			wantErr:  true,
		},
		{
			name:     "invalid name with special characters",
			toolName: "french@greeting",
			wantErr:  true,
		},
		{
			name:     "invalid name too long (65 chars)",
			toolName: strings.Repeat("a", 65),
			wantErr:  true,
		},
		{
			name:     "invalid name with dots",
			toolName: "french.greeting",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tools.ValidateToolName(tt.toolName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToolName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExistingToolNamesValid(t *testing.T) {
	frenchTool := tools.NewFrenchGreetingTool()
	germanTool := tools.NewGermanGreetingTool()

	tests := []struct {
		name     string
		toolName string
	}{
		{
			name:     "french greeting tool name",
			toolName: frenchTool.Definition.Name,
		},
		{
			name:     "german greeting tool name",
			toolName: germanTool.Definition.Name,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tools.ValidateToolName(tt.toolName)
			if err != nil {
				t.Errorf("existing tool name '%s' failed validation: %v", tt.toolName, err)
			}
		})
	}
}