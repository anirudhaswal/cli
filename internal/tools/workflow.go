package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/suprsend-go"
	"gopkg.in/yaml.v3"
)

func triggerWorkflow(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	wfRequestRaw, ok := request.GetArguments()["workflow_request_body"]
	if !ok {
		return mcp.NewToolResultError("Workflow Request Body is required."), nil
	}
	tenantId := request.GetString("tenant_id", "")
	wfRequestBody, ok := wfRequestRaw.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid workflow request body"), nil
	}
	suprsendClient, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	wf := &suprsend.WorkflowTriggerRequest{
		Body:           wfRequestBody,
		IdempotencyKey: "",
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

func newWorkflowTools() []*Tool {
	trigger_suprsend_workflow := &Tool{
		Name:        "workflow.trigger",
		Description: "Enables triggering a specific workflow",
		MCPTool: mcp.NewTool("trigger_suprsend_workflow",
			mcp.WithDescription("Use this tool to trigger a workflow"),
			mcp.WithObject("workflow_request_body",
				mcp.Description("Request body for workflow which will be triggered."),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description("Suprsend workspace to trigger the workflow from."),
				mcp.Required(),
			),
			mcp.WithString("tenant_id",
				mcp.Description("The tenant ID to trigger the workflow for."),
			),
			mcp.WithDestructiveHintAnnotation(true),
		),
		Handler: triggerWorkflow,
	}
	return []*Tool{trigger_suprsend_workflow}
}
func init() {
	for _, w := range newWorkflowTools() {
		RegisterTool(w, "workflow")
	}
}
