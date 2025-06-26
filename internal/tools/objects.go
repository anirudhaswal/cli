package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
	"github.com/suprsend/suprsend-go"
	"gopkg.in/yaml.v3"
)

func getObjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	suprsend_client, err := utils.GetSuprSendWorkspaceClient(workspace)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	obj_identifier := suprsend.ObjectIdentifier{
		ObjectType: object_type,
		Id:         object_id,
	}

	objects_resp, err := suprsend_client.Objects.Get(ctx, obj_identifier)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	yamlobject, err := yaml.Marshal(objects_resp)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(yamlobject)), nil
}

func upsertObjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	action, err := request.RequireString("action")
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

	obj_instance := suprsend_client.Objects.GetEditInstance(obj_identifier)

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

	out, err := utils.HandleObjectAction(ctx, obj_instance, action, key, value, slack_details, identity_provider, obj_identifier, workspace)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(out), nil
}

func newObjectTools() []*Tool {
	get_suprsend_object := &Tool{
		Name:        "objects.get",
		Description: "Enables querying object information",
		MCPTool: mcp.NewTool("get_suprsend_object",
			mcp.WithDescription("Use this tool to get all details about an object"),
			mcp.WithString("object_id",
				mcp.Description("The object_id of the object to get."),
				mcp.Required(),
			),
			mcp.WithString("object_type",
				mcp.Description("The type of object you want to get."),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description("Suprsend workspace to get the object from"),
				mcp.Required(),
			),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		Handler: getObjectHandler,
	}

	upsert_object_handler := &Tool{
		Name:        "objects.upsert",
		Description: "Enables upserting object information",
		MCPTool: mcp.NewTool("upsert_suprsend_object",
			mcp.WithDescription("Use this tool to upsert an object"),
			mcp.WithString("object_id",
				mcp.Description("The object_id of the object to get."),
				mcp.Required(),
			),
			mcp.WithString("object_type",
				mcp.Description("The type of object you want to get."),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description("Suprsend workspace to get the object from."),
				mcp.Required(),
			),
			mcp.WithObject("object_payload",
				mcp.Description("Payload of the request that you want to pass for the object."),
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
		Handler: upsertObjectHandler,
	}

	return []*Tool{get_suprsend_object, upsert_object_handler}
}

func init() {
	for _, t := range newObjectTools() {
		RegisterTool(t, "objects")
	}
}
