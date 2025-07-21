# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Model Context Protocol (MCP) server written in Go using the official MCP Go SDK v0.2.0. The server provides tools that can be used by MCP clients like Claude Desktop.

## Build and Run Commands

```bash
# Build the server
go build

# Run the server (communicates via stdio)
./everyday-mcp-server

# Update dependencies
go mod tidy
```

## Architecture

The MCP server follows a simple architecture:

- **Main Server**: Created using `mcp.NewServer()` with server metadata (name: "everyday", version: "1.0.0")
- **Tool Registration**: Tools are defined with `mcp.Tool` structs and registered using `mcp.AddTool()`
- **Tool Handlers**: Functions that implement the actual tool logic, following the signature `func(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[T]) (*mcp.CallToolResult, error)`
- **Typed Arguments**: Tool arguments use Go structs with JSON and JSONSchema tags for validation
- **Transport**: Uses stdio transport (`mcp.NewStdioTransport()`) for communication with MCP clients

## Tool Development Pattern

When adding new tools:

1. Define argument struct with JSON and JSONSchema tags:
   ```go
   type ToolArgs struct {
       Field string `json:"field" jsonschema:"description of field"`
   }
   ```

2. Create tool handler function:
   ```go
   func toolHandler(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ToolArgs]) (*mcp.CallToolResult, error) {
       // Implementation
       return &mcp.CallToolResult{
           Content: []mcp.Content{
               &mcp.TextContent{Text: "response"},
           },
       }, nil
   }
   ```

3. Register tool in main():
   ```go
   toolDef := &mcp.Tool{
       Name:        "tool-name",
       Description: "Tool description",
   }
   mcp.AddTool(server, toolDef, toolHandler)
   ```

## MCP Client Configuration

To use with Claude Desktop, add to `.mcp.json`:

```json
{
  "mcpServers": {
    "everyday-server": {
      "command": "/path/to/everyday-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```