package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/suprsend-go"
	"gopkg.in/yaml.v3"
)

func triggerWorkflow(ctx context.Context, request mcp.CallToolRequest, workspace, slug string) (*mcp.CallToolResult, error) {
	wfRequestRaw, ok := request.GetArguments()["workflow_request_body"]
	if !ok {
		return mcp.NewToolResultError("Workflow Request Body is required."), nil
	}
	tenantId := request.GetString("tenant_id", "")
	actor := request.GetArguments()["actor"]
	wfRequestBody, ok := wfRequestRaw.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid workflow request body"), nil
	}
	if slug != "" {
		wfRequestBody["workflow"] = slug
	}
	if actor != nil {
		wfRequestBody["actor"] = actor
	}
	suprsendClient, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		log.Error("Error getting workspace client: ", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	wf := &suprsend.WorkflowTriggerRequest{
		Body:           wfRequestBody,
		IdempotencyKey: utils.GenerateUUID(),
		TenantId:       tenantId,
	}

	resp, err := suprsendClient.Workflows.Trigger(wf)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	yamlResp, err := yaml.Marshal(resp)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(yamlResp)), nil
}

func RegisterDynamicWorkflowTools(workspace string) {
	workflows := utils.FetchWorkflowsMcp(workspace)
	if workflows == nil {
		return
	}

	var inputSchema = json.RawMessage(`{
		"type": "object",
		"properties": {
			"workflow_request_body": {
				"type": "object",
				"properties": {
					"recipients": {
						"type": "array",
						"items": {
							"type": "object",
							"additionalProperties": true
						}
					}
				},
				"required": ["recipients"]
			},
			"tenant_id": {
				"type": "string"
			},
			"actor": {
				"type": "object",
				"properties": {
					"distinct_id": {
						"type": "string"
					},
					"is_transient": {
						"type": "boolean"
					}
				}
			}
		}
	}`)

	for _, workflow := range workflows {
		slug := workflow["slug"]
		name := workflow["name"]
		description := workflow["description"]
		if slug == "" {
			continue
		}
		wfTool := &Tool{
			Name:        "workflow." + slug,
			Description: fmt.Sprintf("Trigger workflow: %s - %s", name, description),
			MCPTool: mcp.NewToolWithRawSchema("trigger_workflow_"+slug,
				fmt.Sprintf("Trigger the workflow '%s': %s", name, description),
				inputSchema,
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return triggerWorkflow(ctx, request, workspace, slug)
			},
		}
		RegisterTool(wfTool, "workflow")
	}
}
