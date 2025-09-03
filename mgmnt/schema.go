package mgmnt

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/internal/client"
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
	Type       string                 `json:"type"`
	Defs       map[string]interface{} `json:"$defs"`
	Title      string                 `json:"title"`
	Required   *[]string              `json:"required,omitempty"`
	Properties map[string]interface{} `json:"properties"`
}

type Property struct {
	Type *string `json:"type,omitempty"`
	Ref  *string `json:"$ref,omitempty"`
}

func (c *SS_MgmntClient) ListSchema(workspace string, limit, offset int, mode string) (*ListSchemaResponse, error) {
	if mode != "live" && mode != "draft" {
		return nil, fmt.Errorf("invalid mode: %s. Available modes are: live, draft", mode)
	}
	client := client.NewHTTPClient()
	defer client.Close()

	apiLimit := 50
	allSchemas := []SchemaResponse{}
	currentOffset := offset
	remainingLimit := limit

	for remainingLimit > 0 {
		currentLimit := apiLimit
		if remainingLimit < apiLimit {
			currentLimit = remainingLimit
		}

		url := fmt.Sprintf("%sv1/%s/schema/?limit=%d&offset=%d&mode=%s", c.mgmnt_base_URL, workspace, currentLimit, currentOffset, mode)

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
		if len(schemas.Results) == 0 {
			break
		}

		allSchemas = append(allSchemas, schemas.Results...)
		remainingLimit -= len(schemas.Results)
		currentOffset += len(schemas.Results)
		if len(schemas.Results) < currentLimit {
			break
		}
	}

	return &ListSchemaResponse{
		Results: allSchemas,
		Meta: struct {
			Count  int `json:"count"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Count:  len(allSchemas),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func (c *SS_MgmntClient) GetSchemas(workspace, mode string) (*SchemasResponse, error) {
	if mode != "live" && mode != "draft" {
		return nil, fmt.Errorf("invalid mode: %s, Available modes are: live, draft", mode)
	}

	client := client.NewHTTPClient()
	defer client.Close()

	limit := 50
	offset := 0
	allSchemas := []any{}
	totalCount := 0

	for {
		res, err := client.R().
			SetDebug(c.debug).
			SetHeader("Authorization", "ServiceToken "+c.serviceToken).
			SetResult(&SchemasResponse{}).
			Get(c.mgmnt_base_URL + "v1/" + workspace + "/schema/?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&mode=" + mode)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get schemas: %v\n", err)
			return nil, err
		}

		if res.IsError() {
			fmt.Fprintf(os.Stdout, "Error: Failed to get schemas: %v\n", res.Status())
			return nil, fmt.Errorf("error getting schemas: %s", res.Status())
		}

		schemas := res.Result().(*SchemasResponse)

		if len(schemas.Results) == 0 {
			break
		}

		allSchemas = append(allSchemas, schemas.Results...)
		totalCount += len(schemas.Results)
		offset += limit
	}

	return &SchemasResponse{
		Results: allSchemas,
		Meta: struct {
			Count  int `json:"count"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Count:  totalCount,
			Limit:  limit,
			Offset: 0,
		},
	}, nil
}

func (c *SS_MgmntClient) PushSchema(workspace, schemaSlug string, payload map[string]any) error {
	client := client.NewHTTPClient()
	defer client.Close()
	url := fmt.Sprintf("%sv1/%s/schema/%s/", c.mgmnt_base_URL, workspace, schemaSlug)

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
	return nil
}

func (c *SS_MgmntClient) FinalizeSchema(workspace, slug string, commit bool) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := resty.New()
	defer client.Close()

	urlStr := fmt.Sprintf("%sv1/%s/schema/%s/enable/", c.mgmnt_base_URL, workspace, slug)

	body := map[string]interface{}{
		"is_enabled": commit,
	}

	action := "resetting"
	if commit {
		action = "committing"
	}

	log.Debugf("Finalizing schema (slug: %s) by %s", slug, action)

	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Patch(urlStr)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if res.IsError() {
		if res.StatusCode() == 404 {
			return fmt.Errorf("schema not found: %s", slug)
		}
		return err
	}
	return nil
}
