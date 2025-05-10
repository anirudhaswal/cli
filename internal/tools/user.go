package tools

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"

	log "github.com/sirupsen/logrus"
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

func upsertUserHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinct_id, ok := request.Params.Arguments["distinct_id"].(string)
	if !ok {
		return nil, errors.New("distinct_id must be a string")
	}
	workspace, ok := request.Params.Arguments["workspace"].(string)
	if !ok {
		return nil, errors.New("workspace must be a string")
	}
	action, ok := request.Params.Arguments["action"].(string)
	if !ok {
		return nil, errors.New("action must be a string")
	}
	key, ok := request.Params.Arguments["key"].(string)
	if !ok {
		return nil, errors.New("key must be a string")
	}
	value, ok := request.Params.Arguments["value"].(string)
	if !ok {
		return nil, errors.New("value must be a string")
	}
	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return nil, err
	}
	userInstance := suprsend_client.Users.GetEditInstance(distinct_id)

	switch action {
	case "upsert":
		userInstance.Set(map[string]any{key: value})
		res, err := suprsend_client.Users.AsyncEdit(ctx, userInstance)
		if err != nil {
			return nil, err
		}
		log.Debugf(res.String())
	case "remove":
		userInstance.Remove(map[string]any{key: value})
	case "set":
		userInstance.Set(map[string]any{key: value})
	case "unset":
		userInstance.Unset([]string{key})
	case "append":
		userInstance.Append(map[string]any{key: value})
	case "increment":
		userInstance.Increment(map[string]any{key: value})
	default:
		return nil, errors.New("invalid action")
	}
	return mcp.NewToolResultText("User upserted successfully with distinct_id: " + distinct_id + " and action: " + action + " for key: " + key + " and value: " + value), nil
}

func newUserTools() []*Tool {
	get_suprsend_user := &Tool{
		Name:        "users.get",
		Description: "Enables querying user information",
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

	upsert_suprsend_user := &Tool{
		Name:        "users.upsert",
		Description: "Enables upserting user information",
		MCPTool: mcp.NewTool("upsert_suprsend_user",
			mcp.WithDescription(`Use this tool to upsert a new user or update an existing user's properties.`),
			mcp.WithString("distinct_id",
				mcp.Description(`The distinct_id of the user to get.`),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description(`SuprSend workspace to get the user from.`),
				mcp.Required(),
			),
			mcp.WithString("action",
				mcp.Description(`The action to perform.`),
				mcp.Required(),
				mcp.Enum("upsert", "remove", "set", "unset", "append", "increment"),
			),
			mcp.WithString("key",
				mcp.Description(`The key on which the action is to be performed.`),
				mcp.Required(),
				mcp.Enum("$email", "$phone", "$whatsapp", "$sms", "custom_property", "$timezone", "$preferred_language"),
			),
			mcp.WithString("value",
				mcp.Description(`The value to needs to be added/removed/set/unset/appended/incremented.`),
				mcp.Required(),
			),
		),
		Handler: upsertUserHandler,
	}
	return []*Tool{get_suprsend_user, upsert_suprsend_user}
}

func init() {
	for _, t := range newUserTools() {
		RegisterTool(t, "users")
	}
}
