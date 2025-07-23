package main

import (
	"context"
	"log"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jgsheppa/everyday-mcp-server/pkg/tools"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "everyday",
		Version: "1.1.0",
	}, nil)

	tools.AddFrenchTool(server)
	tools.AddGermanTool(server)

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
