package mgmnt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/internal/client"
)

type EventsResponse struct {
	Meta struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
	Results []any `json:"results"`
}

type ListEventsResponse struct {
	Results []Event `json:"results"`
	Meta    struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
}

type Event struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	PayloadSchema map[string]interface{} `json:"payload_schema"`
	CreatedAt     string                 `json:"created_at"`
}

func (c *SS_MgmntClient) ListEvents(workspace string, limit, offset int) (*ListEventsResponse, error) {
	client := client.NewHTTPClient()
	defer client.Close()

	url := c.mgmnt_base_URL + "v1/" + workspace + "/event/?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
	log.Debugf("Getting Events for workspace: %s, service token: %s", workspace, c.serviceToken)
	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetResult(&ListEventsResponse{}).
		Get(url)
	if err != nil {
		log.Errorf("Error getting events: %s", err)
		return nil, err
	}
	if res.IsError() {
		log.Errorf("Error getting events: %s", res.Status())
		return nil, fmt.Errorf("error getting events: %s", res.Status())
	}
	events := res.Result().(*ListEventsResponse)
	return events, nil
}

func (c *SS_MgmntClient) GetEvents(workspace string) (*EventsResponse, error) {
	client := client.NewHTTPClient()
	defer client.Close()

	limit := 50
	offset := 0
	allEvents := []any{}
	totalCount := 0

	for {
		res, err := client.R().
			SetDebug(c.debug).
			SetHeader("Authorization", "ServiceToken "+c.serviceToken).
			SetResult(&EventsResponse{}).
			Get(c.mgmnt_base_URL + "v1/" + workspace + "/event/?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&has_linked_schema=true")

		if err != nil {
			log.Errorf("Error getting events: %s", err)
			return nil, err
		}
		if res.IsError() {
			log.Errorf("Error getting events: %s", res.Status())
			return nil, fmt.Errorf("error getting events: %s", res.Status())
		}
		events := res.Result().(*EventsResponse)
		if len(events.Results) == 0 {
			break
		}
		allEvents = append(allEvents, events.Results...)
		totalCount += len(events.Results)
		offset += limit
	}
	return &EventsResponse{Results: allEvents}, nil
}

func (c *SS_MgmntClient) PushEvents(workspace string) error {
	client := client.NewHTTPClient()
	defer client.Close()

	filePath := filepath.Join(".", "suprsend", "event", "event_schema_mapping.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("Error reading event schema mapping file: %s", err)
		return err
	}
	var events map[string]any
	if err := json.Unmarshal(data, &events); err != nil {
		log.Errorf("Error parsing event_schema_mapping.json: %s", err)
		return err
	}

	url := c.mgmnt_base_URL + "v1/" + workspace + "/bulk/event/"
	log.Debugf("Pushing events to workspace: %s, service token: %s", workspace, c.serviceToken)
	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetHeader("Content-Type", "application/json").
		SetBody(events).
		Post(url)

	if err != nil {
		log.Errorf("Error pushing event: %s", err)
		return err
	}
	if res.IsError() {
		log.Errorf("Error pushing event: %s", res.Status())
		return fmt.Errorf("error pushing event: %s", res.Status())
	}
	return nil
}
