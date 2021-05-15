package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

type Options struct {
	Repository   *rule.Rules
	BundleClient *bundle.Client
}

type Client struct {
	repository   *rule.Rules
	bundleClient *bundle.Client
}

func NewClient(opts Options) *Client {
	return &Client{
		repository:   opts.Repository,
		bundleClient: opts.BundleClient,
	}
}

func (client *Client) Default(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the opa-bundle-api")
}
