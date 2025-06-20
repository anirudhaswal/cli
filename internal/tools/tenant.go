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
	request_data, ok := request.GetArguments()["tenant_properties"]
	if !ok {
		return mcp.NewToolResultError("tenant_properties is required"), nil
	}
	tenant_properties, ok := request_data.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("invalid tenant_properties"), nil
	}

	tenant_name, tenant_name_ok := tenant_properties["tenant_name"]
	var tenant_name_ptr *string
	if tenant_name_ok {
		tenant_name_str := tenant_name.(string)
		tenant_name_ptr = &tenant_name_str
	}
	tenant_logo, tenant_logo_ok := tenant_properties["logo"]
	var tenant_logo_ptr *string
	if tenant_logo_ok {
		tenant_logo_str := tenant_logo.(string)
		tenant_logo_ptr = &tenant_logo_str
	}
	tenant_timezone, tenant_timezone_ok := tenant_properties["timezone"]
	var tenant_timezone_ptr *string
	if tenant_timezone_ok {
		tenant_timezone_str := tenant_timezone.(string)
		tenant_timezone_ptr = &tenant_timezone_str
	}
	tenant_primary_color, tenant_primary_color_ok := tenant_properties["primary_color"]
	var tenant_primary_color_ptr *string
	if tenant_primary_color_ok {
		tenant_primary_color_str := tenant_primary_color.(string)
		tenant_primary_color_ptr = &tenant_primary_color_str
	}

	tenant_secondary_color, tenant_secondary_color_ok := tenant_properties["secondary_color"]
	var tenant_secondary_color_ptr *string
	if tenant_secondary_color_ok {
		tenant_secondary_color_str := tenant_secondary_color.(string)
		tenant_secondary_color_ptr = &tenant_secondary_color_str
	}
	tenant_tertiary_color, tenant_tertiary_color_ok := tenant_properties["tertiary_color"]
	var tenant_tertiary_color_ptr *string
	if tenant_tertiary_color_ok {
		tenant_tertiary_color_str := tenant_tertiary_color.(string)
		tenant_tertiary_color_ptr = &tenant_tertiary_color_str
	}
	// todo: add blocked_channels
	// blocked_channels, blocked_channels_ok := tenant_properties["blocked_channels"]
	tenant_embedded_preference_url, tenant_embedded_preference_url_ok := tenant_properties["embedded_preference_url"]
	var tenant_embedded_preference_url_ptr *string
	if tenant_embedded_preference_url_ok {
		tenant_embedded_preference_url_str := tenant_embedded_preference_url.(string)
		tenant_embedded_preference_url_ptr = &tenant_embedded_preference_url_str
	}
	tenant_hosted_preference_domain, tenant_hosted_preference_domain_ok := tenant_properties["hosted_preference_domain"]
	var tenant_hosted_preference_domain_ptr *string
	if tenant_hosted_preference_domain_ok {
		tenant_hosted_preference_domain_str := tenant_hosted_preference_domain.(string)
		tenant_hosted_preference_domain_ptr = &tenant_hosted_preference_domain_str
	}
	tenant_social_links, tenant_social_links_ok := tenant_properties["social_links"]
	var tenant_social_links_ptr *suprsend.TenantSocialLinks
	if tenant_social_links_ok {
		social_links_map, ok := tenant_social_links.(map[string]any)
		if ok {
			social_links := &suprsend.TenantSocialLinks{}

			if facebook, ok := social_links_map["facebook"].(string); ok {
				social_links.Facebook = &facebook
			}
			if twitter, ok := social_links_map["twitter"].(string); ok {
				social_links.Twitter = &twitter
			}
			if instagram, ok := social_links_map["instagram"].(string); ok {
				social_links.Instagram = &instagram
			}
			tenant_social_links_ptr = social_links
		}

	}
	tenant_custom_properties, tenant_custom_properties_ok := tenant_properties["custom_properties"]
	var tenant_custom_properties_map map[string]any
	if tenant_custom_properties_ok {
		tenant_custom_properties_map = tenant_custom_properties.(map[string]any)
	}

	tenant_payload := &suprsend.Tenant{
		TenantName:             tenant_name_ptr,
		Logo:                   tenant_logo_ptr,
		Timezone:               tenant_timezone_ptr,
		PrimaryColor:           tenant_primary_color_ptr,
		SecondaryColor:         tenant_secondary_color_ptr,
		TertiaryColor:          tenant_tertiary_color_ptr,
		EmbeddedPreferenceUrl:  tenant_embedded_preference_url_ptr,
		HostedPreferenceDomain: tenant_hosted_preference_domain_ptr,
		SocialLinks:            tenant_social_links_ptr,
		Properties:             tenant_custom_properties_map,
	}

	tenant, err := suprsend_client.Tenants.Upsert(ctx, tenant_id, tenant_payload)
	if err != nil {
		err_str := err.Error()
		if strings.Contains(err_str, "{\"tenant_name\": \"missing value\"}") {
			return mcp.NewToolResultError("tenant_name is required when creating a new tenant try again with a tenant_name"), nil
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
				mcp.Description(`The properties to upsert for the tenant.`),
				mcp.Properties(map[string]any{
					"tenant_name": map[string]any{
						"type":        "string",
						"description": "Name of the tenant",
					},
					"logo": map[string]any{
						"type":        "string",
						"description": "url of the logo of the tenant",
					},
					"timezone": map[string]any{
						"type":        "string",
						"description": "timezone of the tenant",
					},
					"blocked_channels": map[string]any{
						"type":        "array",
						"description": "blocked channels of the tenant",
						"items": map[string]any{
							"type":        "string",
							"description": "blocked channel of the tenant",
						},
					},
					"embedded_preference_url": map[string]any{
						"type":        "string",
						"description": "embedded preference url of the tenant",
					},
					"hosted_preference_domain": map[string]any{
						"type":        "string",
						"description": "hosted preference domain of the tenant",
					},
					"primary_color": map[string]any{
						"type":        "string",
						"description": "primary color of the tenant in hex format",
					},
					"secondary_color": map[string]any{
						"type":        "string",
						"description": "secondary color of the tenant in hex format",
					},
					"tertiary_color": map[string]any{
						"type":        "string",
						"description": "tertiary color of the tenant in hex format",
					},
					"social_links": map[string]any{
						"type":        "object",
						"description": "social links of the tenant",
						"properties": map[string]any{
							"facebook": map[string]any{
								"type":        "string",
								"description": "facebook url of the tenant",
							},
							"twitter": map[string]any{
								"type":        "string",
								"description": "twitter url of the tenant",
							},
							"instagram": map[string]any{
								"type":        "string",
								"description": "instagram url of the tenant",
							},
						},
					},
					"custom_properties": map[string]any{
						"type":        "object",
						"description": "custom properties of the tenant",
					},
				},
				),
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
