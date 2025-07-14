package mgmnt

import (
	"fmt"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"resty.dev/v3"
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

func (c *SS_MgmntClient) GetWorkflows(workspace string, limit int, offset int, mode string) (*WorkflowAPIResponse, error) {
	if mode != "live" && mode != "draft" {
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	client := resty.New()
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

func (c *SS_MgmntClient) PushWorkflow(workspace, slug string, workflow map[string]any) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := resty.New()
	defer client.Close()

	url := fmt.Sprintf("%sv1/%s/workflow/%s/", c.mgmnt_base_URL, workspace, slug)
	fmt.Printf("Request: url %s", url)

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

func (c *SS_MgmntClient) CommitWorkflow(workspace, slug, message string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	client := resty.New()
	defer client.Close()

	urlStr := fmt.Sprintf("%sv1/%s/workflow/%s/commit/", c.mgmnt_base_URL, workspace, slug)

	query := url.Values{}

	if message != "" {
		query.Set("commit_message", message)
	}

	if len(query) > 0 {
		urlStr += "?" + query.Encode()
	}

	log.Debugf("Committing a draft workflow of slug: %s (live)", slug)

	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		Patch(urlStr)

	if err != nil {
		log.Errorf("Error getting workflow: %s", err)
	}

	if res.IsError() {
		log.Errorf("Commit Failed: %s - %s", res.Status(), res.String())
		return fmt.Errorf("push failed: %s", res.Status())
	}

	log.Infof("Successfully committed workflow: %s", slug)
	return nil
}
