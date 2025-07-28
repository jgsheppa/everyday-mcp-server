package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type NotificationArgs struct {
	Message string `json:"message"`
}

func NewNotificationTool() *Config[NotificationArgs] {
	inputSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"message": {
				Type:        "string",
				Description: "The message to send to Google Chat",
			},
		},
		Required: []string{"message"},
	}

	definition := mcp.Tool{
		Name:        "send_notification",
		Description: "Sends a message to a Google Chat space",
		Title:       "Post a message in a Google Chat space",
		InputSchema: inputSchema,
	}

	err := ValidateToolName(definition.Name)
	if err != nil {
		return &Config[NotificationArgs]{
			Error: err,
		}
	}

	return &Config[NotificationArgs]{
		Definition: &definition,
		Call:       SendNotification,
	}
}

func postToGoogleChat(message string) error {
	chatMsg := map[string]string{"text": message}

	jsonData, err := json.Marshal(chatMsg)
	if err != nil {
		return err
	}

	webhookURL := os.Getenv("GOOGLE_CHAT_WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("GOOGLE_CHAT_WEBHOOK_URL environment variable is not set")
	}

	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google Chat API returned status %d", resp.StatusCode)
	}

	return nil
}

func SendNotification(
	_ context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[NotificationArgs],
) (*mcp.CallToolResult, error) {
	message := strings.TrimSpace(params.Arguments.Message)
	if message == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "message argument is empty",
				},
			},
			IsError: true,
		}, nil
	}
	if err := postToGoogleChat(message); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("unable to post message to Google Chat: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "successfully sent a message to Google Chat"},
		},
	}, nil
}

func AddNotificationTool(server *mcp.Server) {
	tool := NewNotificationTool()
	mcp.AddTool(server, tool.Definition, tool.Call)
}
