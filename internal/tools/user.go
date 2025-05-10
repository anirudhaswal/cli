package tools

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
)

func getUserHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinct_id, ok := request.Params.Arguments["distinct_id"].(string)
	if !ok {
		return nil, errors.New("distinct_id must be a string")
	}
	workspace, ok := request.Params.Arguments["workspace"].(string)
	if !ok {
		return nil, errors.New("workspace must be a string")
	}
	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return nil, err
	}
	user, err := suprsend_client.Users.Get(ctx, distinct_id)
	if err != nil {
		return nil, err
	}
	json_user, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(json_user)), nil
}

func newUserTools() []*Tool {
	get_suprsend_user := &Tool{
		Name:        "users.get",
		Description: "Enables querying/updating user information",
		MCPTool: mcp.NewTool("get_suprsend_user",
			mcp.WithDescription(`Use this tool to get all properties for a user in SuprSend. This tool will return a JSON object with all the properties of the user. At top level, it will return the distinct_id, properties (all the custom properties of the user), created_at, updated_at and an array of user channels ($email, push, $sms, $whatsapp, $slack etc.). Eeach object inside will have channel value, status and perma_status (permanent status of the user). If the workspace is not specified. ask the user to provide it before using this tool.`),
			mcp.WithString("distinct_id",
				mcp.Description(`The distinct_id of the user to get.`),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description(`SuprSend workspace to get the user from.`),
				mcp.Required(),
			),
		),
		Handler: getUserHandler,
	}
	return []*Tool{get_suprsend_user}
}

func init() {
	for _, t := range newUserTools() {
		RegisterTool(t, "users")
	}
}
