package client

import (
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
	"resty.dev/v3"
)

func NewHTTPClient() *resty.Client {
	client := resty.New()

	proxyURL := os.Getenv("HTTP_PROXY")
	if proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			log.WithError(err).Error("Invalid HTTP_PROXY")
			return client
		}

		client.SetProxy(parsed.String())
		return client
	}

	return client
}
