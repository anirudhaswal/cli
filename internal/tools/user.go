package tools

import (
	"context"
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
	suprsend "github.com/suprsend/suprsend-go"
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
	value := request.GetString("value", "")

	if utils.RequiresKey(action) && key == "" {
		return mcp.NewToolResultError("key is required for " + action), nil
	}

	if utils.RequiresKey(action) && value == "" {
		return mcp.NewToolResultError("value is required for " + action), nil
	}

	slack_details, err := getSlackDetails(request, action)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	identity_provider := request.GetString("identity_provider", "")
	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	// todo:make everywhere mcp error is returned
	if err != nil {
		return nil, err
	}
	userInstance := suprsend_client.Users.GetEditInstance(distinct_id)

	out, err := utils.HandleUserAction(ctx, userInstance, action, key, value, slack_details, identity_provider, distinct_id, workspace)
	if err != nil {
		return nil, err
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
	user_pref, err := suprsend_client.Users.GetFullPreference(ctx, distinct_id, nil)
	if err != nil {
		return nil, err
	}

	yamluser, err := yaml.Marshal(user_pref)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(yamluser)), nil
}

func getUserCategoryPreference(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinctId, err := request.RequireString("distinct_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	category, err := request.RequireString("category")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	suprsendClient, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return nil, err
	}

	userPref, err := suprsendClient.Users.GetCategoryPreference(ctx, distinctId, category, nil)
	if err != nil {
		return nil, err
	}

	yamlPref, err := yaml.Marshal(userPref)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(yamlPref)), nil
}

func updateUserCategoryPreference(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	distinctId, err := request.RequireString("distinct_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	category, err := request.RequireString("category")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := request.GetArguments()

	rawPayload, ok := args["payload"].(map[string]any)
	if !ok {
		return mcp.NewToolResultError("payload must be an object"), nil
	}

	pref, ok := rawPayload["preference"].(string)
	if !ok {
		return mcp.NewToolResultError("preference must be a string"), nil
	}

	optOutAny, ok := rawPayload["opt_out_channels"]
	if !ok {
		optOutAny = []any{}
	}
	optOutSlice, ok := optOutAny.([]any)
	if !ok {
		return mcp.NewToolResultError("opt_out_channels must be an array"), nil
	}

	optOutChannels := make([]string, 0, len(optOutSlice))
	for _, v := range optOutSlice {
		s, ok := v.(string)
		if !ok {
			return mcp.NewToolResultError("opt_out_channels must be an array of strings"), nil
		}
		optOutChannels = append(optOutChannels, s)
	}

	prefPayload := suprsend.UserUpdateCategoryPreferenceBody{
		Preference:     pref,
		OptOutChannels: optOutChannels,
	}

	workspace, err := request.RequireString("workspace")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	suprsendClient, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return nil, err
	}

	userPref, err := suprsendClient.Users.UpdateCategoryPreference(ctx, distinctId, category, prefPayload, nil)
	if err != nil {
		return nil, err
	}

	yamlPref, err := yaml.Marshal(userPref)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(yamlPref)), nil
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
				mcp.Properties(slackPropertiesSchema),
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

	get_suprsend_category_preference_user := &Tool{
		Name:        "user.get_preferences.category",
		Description: "Enables querying a specific category preference for an user",
		MCPTool: mcp.NewTool("get_suprsend_category_preference_user",
			mcp.WithDescription("Use this tool to query a specific category preference for an user"),
			mcp.WithString("distinct_id",
				mcp.Description("The distinct_id of the user to get."),
				mcp.Required(),
			),
			mcp.WithString("category",
				mcp.Description("The category_slug of an category to get."),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description("SuprSend workspace to run the query from."),
				mcp.Required(),
			),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		Handler: getUserCategoryPreference,
	}

	update_suprsend_category_preference_user := &Tool{
		Name:        "user.update_preferences.category",
		Description: "Enables updating a specific category preference for an user",
		MCPTool: mcp.NewTool("update_suprsend_category_preference_user",
			mcp.WithDescription("Use this tool to update a specific category preference for an user"),
			mcp.WithString("distinct_id",
				mcp.Description("The distinct_id of the user to update."),
				mcp.Required(),
			),
			mcp.WithString("category",
				mcp.Description("category_slug of an category to get."),
				mcp.Required(),
			),
			mcp.WithObject("payload",
				mcp.Description("Payload of an category to update a category preference for an user."),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description("SuprSend workspace to run the query from."),
				mcp.Required(),
			),
			mcp.WithDestructiveHintAnnotation(true),
		),
		Handler: updateUserCategoryPreference,
	}

	return []*Tool{get_suprsend_user, upsert_suprsend_user, get_suprsend_user_preferences, get_suprsend_category_preference_user, update_suprsend_category_preference_user}
}

func init() {
	for _, t := range newUserTools() {
		RegisterTool(t, "users")
	}
}

var slackPropertiesSchema = map[string]any{
	"access_token":               utils.StringSchema("Access token for the Slack workspace"),
	"slack_email":                utils.StringSchema("Email of the user to add to the Slack workspace"),
	"slack_channel_id":           utils.StringSchema("ID of the Slack channel"),
	"slack_user_id":              utils.StringSchema("ID of the Slack user"),
	"slack_incoming_webhook_url": utils.StringSchema("Incoming webhook URL for the Slack"),
}

func getSlackDetails(request mcp.CallToolRequest, action string) (map[string]any, error) {
	if action != "add_slack" && action != "remove_slack" {
		return nil, nil
	}

	slackDetailsRaw, ok := request.GetArguments()["slack_details"]
	if !ok {
		return nil, errors.New("required argument 'slack_details' not found")
	}

	slackDetails, ok := slackDetailsRaw.(map[string]any)
	if !ok {
		return nil, errors.New("invalid slack_details")
	}

	return slackDetails, nil
}
