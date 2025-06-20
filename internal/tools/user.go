package tools

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
	suprsend "github.com/suprsend/suprsend-go"

	log "github.com/sirupsen/logrus"
)

func getUserHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinct_id, err := request.RequireString("distinct_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return nil, err
	}
	user, err := suprsend_client.Users.Get(ctx, distinct_id)
	if err != nil {
		return nil, err
	}

	yamluser, err := yaml.Marshal(user)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(yamluser)), nil
}

func upsertUserHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinct_id, err := request.RequireString("distinct_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	action, err := request.RequireString("action")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	key := request.GetString("key", "")
	if action == "set" || action == "append" || action == "increment" || action == "unset" || action == "remove" {
		if key == "" {
			return mcp.NewToolResultError("key is required for " + action), nil
		}
	}
	value := request.GetString("value", "")
	if action == "set" || action == "append" || action == "increment" {
		if value == "" {
			return mcp.NewToolResultError("value is required for " + action), nil
		}
	}

	var slack_details map[string]any
	if action == "add_slack" || action == "remove_slack" {
		slack_details_raw, ok := request.GetArguments()["slack_details"]
		if !ok {
			return mcp.NewToolResultError("required argument 'slack_details' not found"), nil
		}
		slack_details, ok = slack_details_raw.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid slack_details"), nil
		}
	}

	identity_provider := request.GetString("identity_provider", "")
	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	// todo:make everywhere mcp error is returned
	if err != nil {
		return nil, err
	}
	userInstance := suprsend_client.Users.GetEditInstance(distinct_id)

	var out string

	switch action {
	case "upsert":
		if key != "" && value != "" {
			userInstance.Set(map[string]any{key: value})
			out = "Key set successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
		} else {
			out = "User upserted successfully with distinct_id: " + distinct_id
		}
		res, err := suprsend_client.Users.AsyncEdit(ctx, userInstance)
		if err != nil {
			return nil, err
		}
		log.Debug(res.String())
	case "remove":
		userInstance.Remove(map[string]any{key: value})
		out = "Key removed successfully from user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "set":
		userInstance.Set(map[string]any{key: value})
		out = "Key set successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "unset":
		userInstance.Unset([]string{key})
		out = "Key unset successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "append":
		userInstance.Append(map[string]any{key: value})
		out = "Key appended successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "increment":
		userInstance.Increment(map[string]any{key: value})
		out = "Key incremented successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "add_email":
		userInstance.AddEmail(value)
		out = "Email added successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "set_preferred_language":
		userInstance.SetPreferredLanguage(value)
		out = "Preferred language set successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "set_timezone":
		userInstance.SetTimezone(value)
		out = "Timezone set successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "remove_email":
		userInstance.RemoveEmail(value)
		out = "Email removed successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "add_sms":
		userInstance.AddSms(value)
		out = "SMS added successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "remove_sms":
		userInstance.RemoveSms(value)
		out = "SMS removed successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "add_whatsapp":
		userInstance.AddWhatsapp(value)
		out = "Whatsapp added successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "remove_whatsapp":
		userInstance.RemoveWhatsapp(value)
		out = "Whatsapp removed successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value
	case "add_androidpush":
		userInstance.AddAndroidpush(value, identity_provider)
		out = "Android push added successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value + " and identity_provider: " + identity_provider
	case "remove_androidpush":
		userInstance.RemoveAndroidpush(value, identity_provider)
		out = "Android push removed successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value + " and identity_provider: " + identity_provider
	case "add_iospush":
		userInstance.AddIospush(value, identity_provider)
		out = "iOS push added successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value + " and identity_provider: " + identity_provider
	case "remove_iospush":
		userInstance.RemoveIospush(value, identity_provider)
		out = "iOS push removed successfully for user with distinct_id: " + distinct_id + " for key: " + key + " and value: " + value + " and identity_provider: " + identity_provider
	case "add_slack", "remove_slack":
		slack_incoming_webhook_url, slack_incoming_webhook_url_ok := slack_details["slack_incoming_webhook_url"]
		var payload map[string]any
		var slack_out string
		if slack_incoming_webhook_url_ok {
			payload = map[string]any{"incoming_webhook": map[string]any{"url": slack_incoming_webhook_url}}
			slack_out = "Slack incoming webhook %s successfully for user with distinct_id: " + distinct_id + " and value: " + value
		} else {
			slack_access_token, slack_access_token_ok := slack_details["access_token"]
			slack_email, slack_email_ok := slack_details["slack_email"]
			slack_channel_id, slack_channel_id_ok := slack_details["slack_channel_id"]
			slack_user_id, slack_user_id_ok := slack_details["slack_user_id"]
			if slack_email_ok {
				if !slack_access_token_ok {
					return nil, errors.New("access_token is required when slack_email is provided")
				}
				payload = map[string]any{
					"access_token": slack_access_token,
					"email":        slack_email,
				}
				slack_out = "Slack email %s successfully for user with distinct_id: " + distinct_id + " and value: " + value
			} else if slack_channel_id_ok {
				if !slack_access_token_ok {
					return nil, errors.New("access_token is required when slack_channel_id is provided")
				}
				payload = map[string]any{
					"access_token": slack_access_token,
					"channel_id":   slack_channel_id,
				}
				slack_out = "Slack channel %s successfully for user with distinct_id: " + distinct_id + " and value: " + value
			} else if slack_user_id_ok {
				if !slack_access_token_ok {
					return nil, errors.New("access_token is required when slack_user_id is provided")
				}
				payload = map[string]any{
					"access_token": slack_access_token,
					"user_id":      slack_user_id,
				}
				slack_out = "Slack user id %s successfully for user with distinct_id: " + distinct_id + " and value: " + value
			} else {
				return nil, errors.New("slack_email, slack_channel_id or slack_user_id is required")
			}
		}
		if action == "add_slack" {
			out = fmt.Sprintf(slack_out, "added")
			userInstance.AddSlack(payload)
		} else {
			out = fmt.Sprintf(slack_out, "removed")
			userInstance.RemoveSlack(payload)
		}
	default:
		return nil, errors.New("invalid action")
	}
	_, err = suprsend_client.Users.Edit(ctx, suprsend.UserEditRequest{EditInstance: userInstance})
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(out), nil
}

func getUserPreferencesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinct_id, err := request.RequireString("distinct_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return nil, err
	}
	user_pref, err := suprsend_client.Users.GetUserPreferences(ctx, distinct_id, nil)
	if err != nil {
		return nil, err
	}

	yamluser, err := yaml.Marshal(user_pref)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(yamluser)), nil
}

func newUserTools() []*Tool {
	get_suprsend_user := &Tool{
		Name:        "users.get",
		Description: "Enables querying user information",
		MCPTool: mcp.NewTool("get_suprsend_user",
			mcp.WithDescription(`Use this tool to get all properties for a user in SuprSend. This tool will return a YAML string with all the properties of the user. At top level, it will return the distinct_id, properties (all the custom properties of the user), created_at, updated_at and an array of user channels ($email, push, $sms, $whatsapp, $slack etc.). Eeach object inside will have channel value, status and perma_status (permanent status of the identity). If the workspace is not specified. ask the user to provide it before using this tool.`),
			mcp.WithString("distinct_id",
				mcp.Description(`The distinct_id of the user to get.`),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description(`SuprSend workspace to get the user from.`),
				mcp.Required(),
			),
			mcp.WithReadOnlyHintAnnotation(true),
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
				mcp.Enum(
					"upsert",
					"remove",
					"set",
					"unset",
					"append",
					"increment",
					"add_email",
					"remove_email",
					"add_sms",
					"remove_sms",
					"add_whatsapp",
					"remove_whatsapp",
					"add_androidpush",
					"remove_androidpush",
					"set_preferred_language",
					"set_timezone",
					"add_iospush",
					"remove_iospush",
					"add_slack",
					"remove_slack",
				),
			),
			mcp.WithString("key",
				mcp.Description(`The key on which the action is to be performed. only required for set, append, increment, unset actions.`),
			),
			mcp.WithString("value",
				mcp.Description(`The value to needs to be added/removed/set/unset/appended/incremented.`),
			),
			mcp.WithString("identity_provider",
				mcp.Description(`This is only applicable for add_androidpush, remove_androidpush, add_iospush, remove_iospush actions.`),
			),
			mcp.WithObject("slack_details",
				mcp.Description(`This is only applicable for add_slack and remove_slack actions.`),
				mcp.Properties(map[string]any{
					"access_token": map[string]any{
						"type":        "string",
						"description": "Access token for the Slack workspace",
					},
					"slack_email": map[string]any{
						"type":        "string",
						"description": "Email of the user to add to the Slack workspace",
					},
					"slack_channel_id": map[string]any{
						"type":        "string",
						"description": "ID of the Slack channel",
					},
					"slack_user_id": map[string]any{
						"type":        "string",
						"description": "ID of the Slack user",
					},
					"slack_incoming_webhook_url": map[string]any{
						"type":        "string",
						"description": "Incoming webhook URL for the Slack",
					},
				}),
			),
			mcp.WithDestructiveHintAnnotation(true),
		),
		Handler: upsertUserHandler,
	}
	get_suprsend_user_preferences := &Tool{
		Name:        "users.get_preferences",
		Description: "Enables querying user preferences",
		MCPTool: mcp.NewTool("get_suprsend_user_preferences",
			mcp.WithDescription(`Use this tool to get the preferences for a user in SuprSend.`),
			mcp.WithString("distinct_id",
				mcp.Description(`The distinct_id of the user to get the preferences for.`),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description(`SuprSend workspace to get the user from.`),
				mcp.Required(),
			),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		Handler: getUserPreferencesHandler,
	}
	return []*Tool{get_suprsend_user, upsert_suprsend_user, get_suprsend_user_preferences}
}

func init() {
	for _, t := range newUserTools() {
		RegisterTool(t, "users")
	}
}
