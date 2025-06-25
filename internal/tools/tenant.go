package tools

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/suprsend/cli/internal/utils"
	suprsend "github.com/suprsend/suprsend-go"
	"gopkg.in/yaml.v3"
)

func getTenantHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant_id, err := request.RequireString("tenant_id")
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
	// todo: rename these
	user, err := suprsend_client.Tenants.Get(ctx, tenant_id)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	yamluser, err := yaml.Marshal(user)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(yamluser)), nil
}

func upsertTenantHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tenant_id, err := request.RequireString("tenant_id")
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

	raw_props, ok := request.GetArguments()["tenant_properties"]
	if !ok {
		return mcp.NewToolResultError("tenant_properties is required"), nil
	}

	tenant_properties, ok := raw_props.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid tenant_properties"), nil
	}

	// todo: add blocked_channels
	// blocked_channels, blocked_channels_ok := tenant_properties["blocked_channels"]

	tenant_payload := &suprsend.Tenant{
		TenantName:             utils.GetStringPtr(tenant_properties, "tenant_name"),
		Logo:                   utils.GetStringPtr(tenant_properties, "tenant_logo_ptr"),
		Timezone:               utils.GetStringPtr(tenant_properties, "timezone"),
		PrimaryColor:           utils.GetStringPtr(tenant_properties, "primary_color"),
		SecondaryColor:         utils.GetStringPtr(tenant_properties, "secondary_color"),
		TertiaryColor:          utils.GetStringPtr(tenant_properties, "tertiary_color"),
		EmbeddedPreferenceUrl:  utils.GetStringPtr(tenant_properties, "embedded_preference_url"),
		HostedPreferenceDomain: utils.GetStringPtr(tenant_properties, "hosted_preference_domain"),
		SocialLinks:            buildSocialLinks(utils.GetMap(tenant_properties, "social_links")),
	}

	if custom_tenant_properties, ok := tenant_properties["custom_properties"].(map[string]any); ok {
		tenant_payload.Properties = custom_tenant_properties
	}

	tenant, err := suprsend_client.Tenants.Upsert(ctx, tenant_id, tenant_payload)
	if err != nil {
		err_str := err.Error()
		if strings.Contains(err_str, `{"tenant_name": "missing value"}`) {
			return mcp.NewToolResultError("tenant_name is required when creating a new tenant. Try again with a tenant_name."), nil
		}
		return mcp.NewToolResultError(err_str), nil
	}

	yamluser, err := yaml.Marshal(tenant)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(yamluser)), nil
}

func newTenantTools() []*Tool {
	get_suprsend_tenant := &Tool{
		Name:        "tenants.get",
		Description: "Enables querying tenant information",
		MCPTool: mcp.NewTool("get_suprsend_tenant",
			mcp.WithDescription(`Use this tool to get all properties for a tenant in SuprSend. If the workspace is not specified. ask the user to provide it before using this tool.`),
			mcp.WithString("tenant_id",
				mcp.Description(`The tenant_id of the tenant to get.`),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description(`SuprSend workspace to get the user from.`),
				mcp.Required(),
			),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		Handler: getTenantHandler,
	}

	upsert_suprsend_tenant := &Tool{
		Name:        "tenants.upsert",
		Description: "Enables upserting tenant information",
		MCPTool: mcp.NewTool("upsert_suprsend_tenant",
			mcp.WithDescription(`Use this tool to upsert a new tenant or update an existing tenant's properties.`),
			mcp.WithString("tenant_id",
				mcp.Description(`The tenant_id of the tenant to upsert.`),
				mcp.Required(),
			),
			mcp.WithString("workspace",
				mcp.Description(`SuprSend workspace to get the user from.`),
				mcp.Required(),
			),
			mcp.WithObject("tenant_properties",
				mcp.Description("The properties to upsert for the tenant."),
				mcp.Properties(tenantPropertiesSchema),
			),
			mcp.WithDestructiveHintAnnotation(true),
		),
		Handler: upsertTenantHandler,
	}

	return []*Tool{get_suprsend_tenant, upsert_suprsend_tenant}
}

func init() {
	for _, t := range newTenantTools() {
		RegisterTool(t, "tenants")
	}
}

func buildSocialLinks(m map[string]any) *suprsend.TenantSocialLinks {
	if m == nil {
		return nil
	}
	s := &suprsend.TenantSocialLinks{}
	if v, ok := m["facebook"].(string); ok {
		s.Facebook = &v
	}
	if v, ok := m["twitter"].(string); ok {
		s.Twitter = &v
	}
	if v, ok := m["instagram"].(string); ok {
		s.Instagram = &v
	}
	return s
}

var tenantPropertiesSchema = map[string]any{
	"tenant_name":              utils.StringSchema("Name of the tenant"),
	"logo":                     utils.StringSchema("Url of tenant's logo"),
	"timezone":                 utils.StringSchema("Timezone of the tenant"),
	"blocked_channels":         utils.StringSchema("Blocked channels of the tenant"),
	"embedded_preference_url":  utils.StringSchema("Embedded preference URL"),
	"hosted_preference_domain": utils.StringSchema("Hosted preference domain"),
	"primary_color":            utils.StringSchema("Primary color in hex format"),
	"secondary_color":          utils.StringSchema("Secondary color in hex format"),
	"tertiary_color":           utils.StringSchema("Tertiary color in hex format"),
	"social_links": map[string]any{
		"type":        "object",
		"description": "Social links of the tenant",
		"properties": map[string]any{
			"facebook":  utils.StringSchema("Facebook URL"),
			"twitter":   utils.StringSchema("Twitter URL"),
			"instagram": utils.StringSchema("Instagram URL"),
		},
	},
	"custom_properties": map[string]any{
		"type":        "object",
		"description": "Custom properties of the tenant",
	},
}
