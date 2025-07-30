package mgmnt

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	client := client.NewHTTPClient()
	defer client.Close()

	log.Debugf("Getting workflows for workspace: %s, service token: %s", workspace, c.serviceToken)
	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetResult(&WorkflowAPIResponse{}).
		Get(c.mgmnt_base_URL + "v1/" + workspace + "/workflow/?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&mode=" + mode)
	if err != nil {
		log.Errorf("Error getting workflows: %s", err)
		return nil, err
	}

	if res.IsError() {
		log.Errorf("Error getting workflows: %s", res.Status())
		return nil, fmt.Errorf("error getting workflows: %s", res.Status())
	}

	workflows := res.Result().(*WorkflowAPIResponse)

	return workflows, nil
}

func (c *SS_MgmntClient) GetWorkflows(workspace string, mode string) (*WorkflowsResponse, error) {
	if mode != "live" && mode != "draft" {
		return nil, fmt.Errorf("invalid mode: %s", mode)
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
			fmt.Fprintf(os.Stdout, "No more workflows found, stopping pagination\n")
			break
		}

		allWorkflows = append(allWorkflows, workflows.Results...)
		totalCount += len(workflows.Results)
		offset += limit
	}

	fmt.Fprintf(os.Stdout, "Successfully fetched all %d workflows\n", totalCount)

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

func (c *SS_MgmntClient) PushWorkflow(workspace, slug string, workflow map[string]any) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := client.NewHTTPClient()
	defer client.Close()

	url := fmt.Sprintf("%sv1/%s/workflow/%s/", c.mgmnt_base_URL, workspace, slug)
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

	log.Infof("Successfully pushed workflow: %s", slug)
	return nil
}

func (c *SS_MgmntClient) ChangeStatusWorkflow(workspace, slug string, commit bool) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := client.NewHTTPClient()
	defer client.Close()

	urlStr := fmt.Sprintf("%sv1/%s/workflow/%s/enable/", c.mgmnt_base_URL, workspace, slug)

	body := map[string]interface{}{
		"is_enabled": commit,
	}

	action := "resetting"
	if commit {
		action = "committing"
	}

	log.Debugf("Finalizing workflow (slug: %s) by %s", slug, action)

	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Patch(urlStr)

	if err != nil {
		log.Errorf("Error %s workflow (slug: %s): %s", action, slug, err)
		return err
	}

	if res.IsError() {
		log.Errorf("%s failed for workflow (slug: %s): %s - %s", action, slug, res.Status(), res.String())
		return fmt.Errorf("%s failed: %s - %s", action, res.Status(), res.String())
	}

	log.Infof("Successfully %s workflow: %s", strings.TrimSuffix(action, "ing")+"ed", slug)
	return nil
}
