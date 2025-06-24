package mgmnt

import (
	"fmt"
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

func (c *SS_MgmntClient) GetWorkflows(workspace string, limit int, offset int, mode string) ([]Workflow, error) {
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

	return workflows.Results, nil
}
