package mgmnt

import (
	"fmt"
	"log"

	"resty.dev/v3"
)

type SchemasResponse struct {
	Meta struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
	Results []any `json:"results"`
}

type ListSchemaResponse struct {
	Meta struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
	Results []SchemaResponse `json:"results"`
}

type SchemaResponse struct {
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsEnabled   bool       `json:"is_enabled"`
	JSONSchema  JSONSchema `json:"json_schema"`
}

type SchemaPayload struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsEnabled   string     `json:"is_enabled"`
	JSONSchema  JSONSchema `json:"json_schema"`
}

type JSONSchema struct {
	Type       string              `json:"type"`
	Title      string              `json:"title"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type string `json:"type"`
}

func (c *SS_MgmntClient) ListSchema(workspace string) (*ListSchemaResponse, error) {
	client := resty.New()
	defer client.Close()

	url := fmt.Sprintf("%s/v1/%s/schema/", c.baseURL, workspace)

	resp, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetResult(&ListSchemaResponse{}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	schemas := resp.Result().(*ListSchemaResponse)
	return schemas, nil
}

func (c *SS_MgmntClient) GetSchemas(workspace string) (*SchemasResponse, error) {
	client := resty.New()
	defer client.Close()

	url := fmt.Sprintf("%s/v1/%s/schema/", c.baseURL, workspace)

	resp, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetResult(&SchemasResponse{}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	schemas := resp.Result().(*SchemasResponse)
	return schemas, nil
}

func (c *SS_MgmntClient) GetSchema(workspace, schemaSlug string) (*SchemaResponse, error) {
	client := resty.New()
	defer client.Close()

	url := fmt.Sprintf("%s/v1/%s/schema/%s/", c.baseURL, workspace, schemaSlug)

	resp, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetResult(&SchemaResponse{}).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.Status())
	}

	schema := resp.Result().(*SchemaResponse)
	return schema, nil
}

func (c *SS_MgmntClient) PushSchema(workspace, schemaSlug string, payload map[string]any) error {
	client := resty.New()
	defer client.Close()

	url := fmt.Sprintf("%s/v1/%s/schema/%s/", c.baseURL, workspace, schemaSlug)

	resp, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("error response: %s", resp.Status())
	}

	log.Printf("Pushed schema to %s", workspace)
	return nil
}
