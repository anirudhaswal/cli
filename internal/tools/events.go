package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
)

func triggerEvent(ctx context.Context, request mcp.CallToolRequest, workspace, name string) (*mcp.CallToolResult, error) {
	eventRequestRaw, ok := request.GetArguments()["payload_schema"]
	if !ok {
		return mcp.NewToolResultError("payload schema is required"), nil
	}

	eventRequestBody, ok := eventRequestRaw.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("payload schema is not a valid object"), nil
	}

	schema := eventRequestBody["schema"].(string)
	versionNo := eventRequestBody["version_no"].(string)
	suprsendClient, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get suprsend client: %v", err)), nil
	}

	// trigger the event thru go sdk

	return mcp.NewToolResultText("Event triggered successfully"), nil
}

func RegisterDynamicEventsTools(workspace string) error {
	events := utils.FetchEventsMcp(workspace)
	if len(events) == 0 {
		return fmt.Errorf("No workflows present in %s workspace", workspace)
	}

	var inputSchema = json.RawMessage(`{
		"type": "object",
		"properties": {
			"payload_schema": {
				"type": "object",
				"properties": {
					"schema": {
						"type": "string"
					},
					"version_no": {
						"type": "string"
					}
				}
			}
		}
	}`)

	for _, event := range events {
		name := event["name"]
		description := event["description"]
		if name == "" {
			continue
		}
		eventTool := &Tool{
			Name:        "trigger_event." + name,
			Description: fmt.Sprintf("Trigger event %s event - %s", name, description),
			MCPTool: mcp.NewToolWithRawSchema("trigger_event_"+name,
				fmt.Sprintf("Trigger %s event - %s", name, description),
				inputSchema,
			),
		}
		RegisterTool(eventTool, "event")
	}
	return nil
}
