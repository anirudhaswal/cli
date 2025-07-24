package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/suprsend-go"
	"gopkg.in/yaml.v3"
)

func getObjectSubscriptionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	object_id, err := request.RequireString("object_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	object_type, err := request.RequireString("object_type")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit_subscriptions, err := request.RequireInt("limit")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cursor_list_api_opts := suprsend.CursorListApiOptions{
		Limit: limit_subscriptions,
	}
	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	obj_identifier := suprsend.ObjectIdentifier{
		ObjectType: object_type,
		Id:         object_id,
	}

	obj_subs_resp, err := suprsend_client.Objects.GetSubscriptions(ctx, obj_identifier, &cursor_list_api_opts)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	yamlobject, err := yaml.Marshal(obj_subs_resp)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(yamlobject)), nil
}

func addObjectSubscriptionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	object_id, err := request.RequireString("object_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	object_type, err := request.RequireString("object_type")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	recipients, ok := request.GetArguments()["recipients"]
	if !ok {
		return mcp.NewToolResultError("Recipients is a required property"), nil
	}

	props, ok := request.GetArguments()["properties"]
	if !ok {
		return mcp.NewToolResultError("Properties isn't passed as an object"), nil
	}

	parent_obj_props, ok := request.GetArguments()["parent_object_properties"]
	if !ok {
		return mcp.NewToolResultError("Parent object_properties isn't passed properly"), nil
	}

	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	obj_identifier := suprsend.ObjectIdentifier{
		ObjectType: object_type,
		Id:         object_id,
	}

	payload := map[string]any{
		"recipients":               recipients,
		"parent_object_properties": parent_obj_props,
		"properties":               props,
	}

	obj_subs_resp, err := suprsend_client.Objects.CreateSubscriptions(ctx, obj_identifier, payload)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	yamlobject, err := yaml.Marshal(obj_subs_resp)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(yamlobject)), nil
}
