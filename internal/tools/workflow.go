package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/internal/commands/schema"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/suprsend-go"
)

func triggerWorkflow(_ context.Context, request mcp.CallToolRequest, workspace, slug string) (*mcp.CallToolResult, error) {
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

	_, ok = wfRequestBody["schema"]
	if !ok {
		responseData := map[string]any{
			"response": resp,
			"warning":  fmt.Sprintf("Schema is not present for %s workflow", slug),
		}
		jsonData, err := json.MarshalIndent(responseData, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultStructured(responseData, string(jsonData)), nil
	}
	responseData := map[string]any{
		"response": resp,
	}

	jsonData, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultStructured(responseData, string(jsonData)), nil
}

func RegisterDynamicWorkflowTools(workspace, workflowsFlag string) error {
	workflows := utils.FetchWorkflowsMcp(workspace, workflowsFlag)
	if len(workflows) == 0 {
		return fmt.Errorf("No workflows present in %s workspace", workspace)
	}

	var patchSchema = json.RawMessage(`{
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
		slug := workflow.Slug
		if slug == "" {
			continue
		}
		name := workflow.Name
		description := workflow.Description
		mgmntClient := utils.GetSuprSendMgmntClient()
		payloadSchema, err := mgmntClient.GetSchema(workspace, workflow.PayloadSchema.Schema, workflow.PayloadSchema.Version)
		if err != nil {
			return fmt.Errorf("failed to get schema: %s", err)
		}
		inputSchema, err := json.Marshal(payloadSchema.JSONSchema)
		if err != nil {
			return fmt.Errorf("failed to marshal schema: %s", err)
		}
		mergedSchema, err := schema.MergeAndValidate(string(inputSchema), []byte(patchSchema))
		if err != nil {
			log.Errorf("failed to create workflow schema for wf %s: %s", name, err)
			continue
		}
		if name == "" {
			continue
		} else {
			name = strings.ReplaceAll(name, " ", "_")
			name = strings.ToLower(name)
		}
		if description == "" {
			description = fmt.Sprintf("Use this tool to trigger workflow with name: \"%s\"", name)
		} else {
			description = fmt.Sprintf("Use this tool to trigger workflow with name: \"%s\" with description:\"%s\"", name, description)
		}
		cleanSlug := strings.ReplaceAll(slug, "-", "_")
		wfTool := &Tool{
			Name:        "trigger_" + cleanSlug + "_workflow",
			Description: fmt.Sprintf("Trigger workflow: %s - %s", name, description),
			MCPTool: mcp.NewToolWithRawSchema("trigger_"+cleanSlug+"_workflow",
				description,
				mergedSchema,
			),
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return triggerWorkflow(ctx, request, workspace, slug)
			},
		}
		RegisterWorkflow(wfTool)
	}
	return nil
}
