package tools

import (
	"context"
	"fmt"
	"time"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type YearsSinceArgs struct {
	Year int `json:"year" jsonschema:"the year to calculate years since"`
}

func YearsSince(
	_ context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[YearsSinceArgs],
) (*mcp.CallToolResult, error) {
	targetYear := params.Arguments.Year
	currentYear := time.Now().Year()
	
	if targetYear > currentYear {
		return nil, fmt.Errorf("year %d is in the future", targetYear)
	}
	
	yearsPassed := currentYear - targetYear
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%d years have passed since %d (current year: %d)", yearsPassed, targetYear, currentYear),
			},
		},
	}, nil
}

func NewYearsSinceTool() *Config[YearsSinceArgs] {
	definition := mcp.Tool{
		Name:        "years_since",
		Description: "Calculates how many years have passed since a specific year",
		Title:       "Years Since",
	}

	err := ValidateToolName(definition.Name)
	if err != nil {
		return &Config[YearsSinceArgs]{
			Error: err,
		}
	}

	return &Config[YearsSinceArgs]{
		Definition: &definition,
		Call:       YearsSince,
	}
}

func AddYearsSinceTool(server *mcp.Server) {
	tool := NewYearsSinceTool()
	mcp.AddTool(server, tool.Definition, tool.Call)
}