package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	suprsend "github.com/suprsend/suprsend-go"
)

func McpStringField(name, description string, required bool, enumValues ...string) mcp.ToolOption {
	opts := []mcp.PropertyOption{
		mcp.Description(description),
	}
	if required {
		opts = append(opts, mcp.Required())
	}
	if len(enumValues) > 0 {
		opts = append(opts, mcp.Enum(enumValues...))
	}
	return mcp.WithString(name, opts...)
}

func GetStringPtr(m map[string]any, key string) *string {
	if val, ok := m[key].(string); ok {
		return &val
	}
	return nil
}

func GetMap(m map[string]any, key string) map[string]any {
	if val, ok := m[key].(map[string]any); ok {
		return val
	}
	return nil
}

func StringSchema(desc string) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": desc,
	}
}

func ArraySchema(desc string) map[string]any {
	return map[string]any{
		"type":        "array",
		"description": desc,
		"items": map[string]any{
			"type": "string",
		},
	}
}

func RequiresKey(action string) bool {
	actions := map[string]bool{
		"set": true, "append": true, "increment": true, "unset": true, "remove": true,
	}
	return actions[action]
}

func HandleAction(ctx context.Context, userInstance suprsend.UserEdit, action, key, value string, slack_details map[string]interface{}, identity_provider, distinct_id string, workspace string) (string, error) {
	var err error
	var out string

	suprsend_client, err := GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return "", err
	}

	switch action {
	case "upsert":
		if key != "" && value != "" {
			userInstance.Set(map[string]interface{}{key: value})
			out = fmt.Sprintf("Key set successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
		} else {
			out = "User upserted successfully with distinct_id: " + distinct_id
		}
		_, err = suprsend_client.Users.AsyncEdit(ctx, userInstance)
	case "remove":
		userInstance.Remove(map[string]interface{}{key: value})
		out = fmt.Sprintf("Key removed successfully from user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "set":
		userInstance.Set(map[string]interface{}{key: value})
		out = fmt.Sprintf("Key set successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "unset":
		userInstance.Unset([]string{key})
		out = fmt.Sprintf("Key unset successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "append":
		userInstance.Append(map[string]interface{}{key: value})
		out = fmt.Sprintf("Key appended successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "increment":
		userInstance.Increment(map[string]interface{}{key: value})
		out = fmt.Sprintf("Key incremented successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "add_email":
		userInstance.AddEmail(value)
		out = fmt.Sprintf("Email added successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "set_preferred_language":
		userInstance.SetPreferredLanguage(value)
		out = fmt.Sprintf("Preferred language set successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "set_timezone":
		userInstance.SetTimezone(value)
		out = fmt.Sprintf("Timezone set successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "remove_email":
		userInstance.RemoveEmail(value)
		out = fmt.Sprintf("Email removed successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "add_sms":
		userInstance.AddSms(value)
		out = fmt.Sprintf("SMS added successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "remove_sms":
		userInstance.RemoveSms(value)
		out = fmt.Sprintf("SMS removed successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "add_whatsapp":
		userInstance.AddWhatsapp(value)
		out = fmt.Sprintf("Whatsapp added successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "remove_whatsapp":
		userInstance.RemoveWhatsapp(value)
		out = fmt.Sprintf("Whatsapp removed successfully for user with distinct_id: %s for key: %s and value: %s", distinct_id, key, value)
	case "add_androidpush":
		userInstance.AddAndroidpush(value, identity_provider)
		out = fmt.Sprintf("Android push added successfully for user with distinct_id: %s for key: %s and value: %s and identity_provider: %s", distinct_id, key, value, identity_provider)
	case "remove_androidpush":
		userInstance.RemoveAndroidpush(value, identity_provider)
		out = fmt.Sprintf("Android push removed successfully for user with distinct_id: %s for key: %s and value: %s and identity_provider: %s", distinct_id, key, value, identity_provider)
	case "add_iospush":
		userInstance.AddIospush(value, identity_provider)
		out = fmt.Sprintf("iOS push added successfully for user with distinct_id: %s for key: %s and value: %s and identity_provider: %s", distinct_id, key, value, identity_provider)
	case "remove_iospush":
		userInstance.RemoveIospush(value, identity_provider)
		out = fmt.Sprintf("iOS push removed successfully for user with distinct_id: %s for key: %s and value: %s and identity_provider: %s", distinct_id, key, value, identity_provider)
	case "add_slack", "remove_slack":
		payload, slackOut, err := prepareSlackPayload(slack_details, action)
		if err != nil {
			return "", err
		}
		if action == "add_slack" {
			userInstance.AddSlack(payload)
		} else {
			userInstance.RemoveSlack(payload)
		}
		out = fmt.Sprintf(slackOut, distinct_id, value)
	}

	return out, err
}

func prepareSlackPayload(slackDetails map[string]interface{}, action string) (map[string]interface{}, string, error) {
	var payload map[string]interface{}
	var slackOut string

	if url, ok := slackDetails["slack_incoming_webhook_url"]; ok {
		payload = map[string]interface{}{"incoming_webhook": map[string]interface{}{"url": url}}
		slackOut = "Slack incoming webhook %s successfully for user with distinct_id: %s and value: %s"
	} else {
		accessToken, accessTokenOk := slackDetails["access_token"]
		slackEmail, slackEmailOk := slackDetails["slack_email"]
		slackChannelID, slackChannelIDOk := slackDetails["slack_channel_id"]
		slackUserID, slackUserIDOk := slackDetails["slack_user_id"]

		if slackEmailOk {
			if !accessTokenOk {
				return nil, "", errors.New("access_token is required when slack_email is provided")
			}
			payload = map[string]interface{}{
				"access_token": accessToken,
				"email":        slackEmail,
			}
			slackOut = "Slack email %s successfully for user with distinct_id: %s and value: %s"
		} else if slackChannelIDOk {
			if !accessTokenOk {
				return nil, "", errors.New("access_token is required when slack_channel_id is provided")
			}
			payload = map[string]interface{}{
				"access_token": accessToken,
				"channel_id":   slackChannelID,
			}
			slackOut = "Slack channel %s successfully for user with distinct_id: %s and value: %s"
		} else if slackUserIDOk {
			if !accessTokenOk {
				return nil, "", errors.New("access_token is required when slack_user_id is provided")
			}
			payload = map[string]interface{}{
				"access_token": accessToken,
				"user_id":      slackUserID,
			}
			slackOut = "Slack user id %s successfully for user with distinct_id: %s and value: %s"
		} else {
			return nil, "", errors.New("slack_email, slack_channel_id or slack_user_id is required")
		}
	}

	return payload, slackOut, nil
}
