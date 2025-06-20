package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Type        string
	Name        string
	Description string
	MCPTool     mcp.Tool
	Handler     func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

var toolRegistry []*Tool

func RegisterTool(t *Tool, toolType string) {
	t.Type = toolType
	toolRegistry = append(toolRegistry, t)
}

func GetAllTools() []*Tool {
	return toolRegistry
}
