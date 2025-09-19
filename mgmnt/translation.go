package mgmnt

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/suprsend/cli/internal/client"
)

type TranslationItem struct {
	Name      string `json:"name"`
	Locale    string `json:"locale"`
	FileName  string `json:"filename"`
	VersionNo int    `json:"version_no"`
	Status    string `json:"status"`
	Action    string `json:"action"`
}

type ListTranslation struct {
	Results []TranslationItem `json:"results"`
	Meta    struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
}

type TranslationResponse struct {
	Results []any `json:"results"`
	Meta    struct {
		Count  int `json:"count"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
}

func (c *SS_MgmntClient) ListTranslations(workspace, mode string) (*ListTranslation, error) {
	client := client.NewHTTPClient()
	defer client.Close()

	url := c.mgmnt_base_URL + "v1/" + workspace + "/translation/?mode=" + mode
	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetResult(&ListTranslation{}).
		Get(url)

	if err != nil {
		log.Errorf("Error getting translations: %s", err)
		return nil, err
	}
	if res.IsError() {
		log.Errorf("Error getting Translation: %s", res.Status())
		return nil, err
	}

	translations := res.Result().(*ListTranslation)

	return translations, nil
}

func (c *SS_MgmntClient) GetTranslations(workspace, mode string) (*TranslationResponse, error) {
	if mode != "live" && mode != "draft" {
		log.Errorf("%s: invalid mode. Available modes are: draft, live", mode)
		return nil, nil
	}
	client := client.NewHTTPClient()
	defer client.Close()

	limit := 10
	offset := 0
	allTranslations := []any{}
	totalCount := 0

	for {
		res, err := client.R().
			SetDebug(c.debug).
			SetHeader("Authorization", "ServiceToken "+c.serviceToken).
			SetResult(&TranslationResponse{}).
			Get(c.mgmnt_base_URL + "v1/" + workspace + "/translation/?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&mode=" + mode)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to get translations: %v\n", err)
			return nil, err
		}

		if res.IsError() {
			fmt.Fprintf(os.Stdout, "Error: Failed to get translations: %v\n", res.Status())
			return nil, fmt.Errorf("error getting translations: %s", res.Status())
		}

		translations := res.Result().(*TranslationResponse)

		if len(translations.Results) == 0 {
			break
		}

		allTranslations = append(allTranslations, translations.Results...)
		totalCount += len(translations.Results)
		offset += limit
	}

	return &TranslationResponse{
		Results: allTranslations,
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

func (c *SS_MgmntClient) PushTranslation(workspace string, translations []map[string]any) error {
	client := client.NewHTTPClient()
	defer client.Close()
	url := c.mgmnt_base_URL + "v1/" + workspace + "/translation/"
	body := map[string]any{
		"translations": translations,
	}
	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetBody(body).
		Post(url)
	if err != nil {
		log.Errorf("Error pushing translations: %s", err)
		return err
	}
	if res.IsError() {
		log.Errorf("Error pushing translations: %s", res.Status())
		return fmt.Errorf("error pushing translations: %s", res.Status())
	}
	return nil
}

func (c *SS_MgmntClient) FinalizeTranslation(workspace, commitMessage string) error {
	client := client.NewHTTPClient()
	defer client.Close()
	url := c.mgmnt_base_URL + "v1/" + workspace + "/translation/commit/?commit_message=" + commitMessage
	res, err := client.R().
		SetDebug(c.debug).
		SetHeader("Authorization", "ServiceToken "+c.serviceToken).
		SetBody("{}").
		Patch(url)
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("Error committing translation: %s", res.Status())
	}
	return nil
}
