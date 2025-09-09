package mgmnt

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/internal/client"
)

type Workflow struct {
	Slug      string   `json:"slug"`
	IsEnabled bool     `json:"is_enabled"`
	Status    string   `json:"status"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
}

type WorkflowAPIResponse struct {
	Results []Workflow `json:"results"`
	Meta    struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
}

type WorkflowsResponse struct {
	Results []any `json:"results"`
	Meta    struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
}

func (c *SS_MgmntClient) ListWorkflows(workspace string, limit int, offset int, mode string) (*WorkflowAPIResponse, error) {
	if mode != "live" && mode != "draft" {
		return nil, fmt.Errorf("invalid mode: %s. Available modes are: live, draft", mode)
	}

	client := client.NewHTTPClient()
	defer client.Close()

	apiLimit := 50
	allWorkflows := []Workflow{}
	currentOffset := offset
	remainingLimit := limit
	for remainingLimit > 0 {
		currentLimit := apiLimit
		if remainingLimit < apiLimit {
			currentLimit = remainingLimit
		}

		log.Debugf("Getting workflows for workspace: %s, service token: %s, limit: %d, offset: %d", workspace, c.serviceToken, currentLimit, currentOffset)
		res, err := client.R().
			SetDebug(c.debug).
			SetHeader("Authorization", "ServiceToken "+c.serviceToken).
			SetResult(&WorkflowAPIResponse{}).
			Get(c.mgmnt_base_URL + "v1/" + workspace + "/workflow/?limit=" + strconv.Itoa(currentLimit) + "&offset=" + strconv.Itoa(currentOffset) + "&mode=" + mode)
		if err != nil {
			log.Errorf("Error getting workflows: %s", err)
			return nil, err
		}
		if res.IsError() {
			log.Errorf("Error getting workflows: %s", res.Status())
			return nil, fmt.Errorf("error getting workflows: %s", res.Status())
		}

		workflows := res.Result().(*WorkflowAPIResponse)
		if len(workflows.Results) == 0 {
			break
		}

		allWorkflows = append(allWorkflows, workflows.Results...)
		remainingLimit -= len(workflows.Results)
		currentOffset += len(workflows.Results)
		if len(workflows.Results) < currentLimit {
			break
		}
	}

	return &WorkflowAPIResponse{
		Results: allWorkflows,
		Meta: struct {
			Count  int `json:"count"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Count:  len(allWorkflows),
			Limit:  limit,
			Offset: currentOffset,
		},
	}, nil
}

func (c *SS_MgmntClient) GetWorkflows(workspace, mode string) (*WorkflowsResponse, error) {
	if mode != "live" && mode != "draft" {
		return nil, fmt.Errorf("invalid mode: %s, Available modes are: live, draft", mode)
	}

	client := client.NewHTTPClient()
	defer client.Close()

	limit := 50
	offset := 0
	allWorkflows := []any{}
	totalCount := 0

	for {
		res, err := client.R().
			SetDebug(c.debug).
			SetHeader("Authorization", "ServiceToken "+c.serviceToken).
			SetResult(&WorkflowsResponse{}).
			Get(c.mgmnt_base_URL + "v1/" + workspace + "/workflow/?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&mode=" + mode)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get workflows: %v\n", err)
			return nil, err
		}

		if res.IsError() {
			fmt.Fprintf(os.Stdout, "Error: Failed to get workflows: %v\n", res.Status())
			return nil, fmt.Errorf("error getting workflows: %s", res.Status())
		}

		workflows := res.Result().(*WorkflowsResponse)

		if len(workflows.Results) == 0 {
			break
		}

		allWorkflows = append(allWorkflows, workflows.Results...)
		totalCount += len(workflows.Results)
		offset += limit
	}

	return &WorkflowsResponse{
		Results: allWorkflows,
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

func (c *SS_MgmntClient) PushWorkflow(workspace, slug string, workflow map[string]any, commit bool, commitMessage string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := client.NewHTTPClient()
	defer client.Close()

	url := fmt.Sprintf("%sv1/%s/workflow/%s/?commit=%t&commit_message=%s", c.mgmnt_base_URL, workspace, slug, commit, commitMessage)
	log.Debugf("Pushing workflow to: %s", url)

	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetBody(workflow).
		Post(url)

	if err != nil {
		log.Errorf("Error pushing workflow: %s", err)
		return err
	}
	if res.IsError() {
		log.Errorf("Push failed: %s - %s", res.Status(), res.String())
		return fmt.Errorf("push failed: %s", res.Status())
	}
	return nil
}

func (c *SS_MgmntClient) ChangeStatusWorkflow(workspace, slug string, enabled bool) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := client.NewHTTPClient()
	defer client.Close()

	urlStr := fmt.Sprintf("%sv1/%s/workflow/%s/enable/", c.mgmnt_base_URL, workspace, slug)

	body := map[string]interface{}{
		"is_enabled": enabled,
	}

	action := "disabling"
	if enabled {
		action = "enabling"
	}

	log.Debugf("Finalizing workflow (slug: %s) by %s", slug, action)

	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Patch(urlStr)

	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}

	if res.IsError() {
		if res.StatusCode() == 404 {
			return fmt.Errorf("workflow not found: %s", slug)
		}
		return fmt.Errorf("%s failed: %s", action, res.Status())
	}

	return nil
}
