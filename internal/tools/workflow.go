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
	wfRequestRaw := request.GetArguments()
	tenantId := request.GetString("tenant_id", "")
	suprsendClient, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		log.Error("Error getting workspace client: ", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	// add workflow slug to the request body
	wfRequestBody := map[string]any{}
	wfRequestBody["workflow"] = slug
	// if tenant id is present, add it to the request body
	if tenantId != "" {
		wfRequestBody["tenant_id"] = tenantId
	}
	// actor: object { distinct_id: string }
	if actorVal, ok := wfRequestRaw["actor"]; ok {
		if actorMap, ok := actorVal.(map[string]any); ok {
			if _, ok := actorMap["distinct_id"].(string); !ok {
				return mcp.NewToolResultError("invalid 'actor': 'distinct_id' must be a string"), nil
			}
			wfRequestBody["actor"] = actorMap
		} else {
			return mcp.NewToolResultError("invalid 'actor': must be an object"), nil
		}
	} else if v, ok := wfRequestRaw["actor.distinct_id"]; ok {
		// Back-compat support for dotted key
		if s, ok := v.(string); ok {
			wfRequestBody["actor"] = map[string]any{"distinct_id": s}
		} else {
			return mcp.NewToolResultError("invalid 'actor.distinct_id': must be a string"), nil
		}
	}
	// if recipients is present, add it to the request body
	if recipients, ok := wfRequestRaw["recipients"]; ok {
		wfRequestBody["recipients"] = recipients
	}
	// Add data to the request body
	wfRequestBody["data"] = wfRequestRaw["data"]
	wf := &suprsend.WorkflowTriggerRequest{
		Body:           wfRequestBody,
		IdempotencyKey: utils.GenerateUUID(),
		TenantId:       tenantId,
	}

	resp, err := suprsendClient.Workflows.Trigger(wf)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
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
		return fmt.Errorf("no workflows present in %s workspace", workspace)
	}

	patchSchema := json.RawMessage(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "TenantActorRecipientsSchema",
  "type": "object",
  "properties": {
    "tenant_id": {
      "type": "string",
      "default": "default",
      "description": "Unique identifier for the tenant. Defaults to 'default' if not provided."
    },
    "actor": {
      "type": "object",
      "description": "Represents the actor who initiated the action, identified by a distinct_id.",
      "properties": {
        "distinct_id": {
          "type": "string",
          "description": "Unique identifier of the actor (e.g., user ID).",
          "example": "kfdfdj"
        }
      },
      "required": ["distinct_id"],
      "additionalProperties": false
    },
    "recipients": {
      "type": "array",
      "description": "List of recipient distinct_ids. Each item in this array is a string representing a recipient's unique identifier.",
      "items": {
        "type": "string",
        "example": "kfdfdj"
      },
      "minItems": 1
    }
  },
  "required": ["recipients"],
  "additionalProperties": false
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
		mergedSchema, err := schema.MergeUnderDataAndValidate(string(patchSchema), string(inputSchema))
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
