package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/logs"
	"github.com/xenitab/opa-bundle-api/pkg/replay"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

type Options struct {
	RuleClient   *rule.Client
	BundleClient *bundle.Client
	LogsClient   *logs.Client
	ReplayClient *replay.Client
}

type Client struct {
	ruleClient   *rule.Client
	bundleClient *bundle.Client
	logsClient   *logs.Client
	replayClient *replay.Client
}

func NewClient(opts Options) *Client {
	return &Client{
		ruleClient:   opts.RuleClient,
		bundleClient: opts.BundleClient,
		logsClient:   opts.LogsClient,
		replayClient: opts.ReplayClient,
	}
}

func (client *Client) Default(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the opa-bundle-api")
}
