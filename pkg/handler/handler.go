package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/logs"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

type Options struct {
	Repository   *rule.Rules
	BundleClient *bundle.Client
	LogsClient   *logs.Logs
}

type Client struct {
	repository   *rule.Rules
	bundleClient *bundle.Client
	logsClient   *logs.Logs
}

func NewClient(opts Options) *Client {
	return &Client{
		repository:   opts.Repository,
		bundleClient: opts.BundleClient,
		logsClient:   opts.LogsClient,
	}
}

func (client *Client) Default(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the opa-bundle-api")
}
