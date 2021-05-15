package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/util"
)

func (client *Client) GetBundle(c echo.Context) error {
	data, err := client.repository.GetAllJSON()
	if err != nil {
		return err
	}

	dataBytes := []byte(data)
	revision, err := util.BytesToHash(dataBytes)

	req := c.Request()
	headers := req.Header
	headerIfNoneMatch := headers.Get("If-None-Match")

	if headerIfNoneMatch == revision {
		return c.NoContent(http.StatusNotModified)
	}

	bundleClient := bundle.NewClient()
	archive, err := bundleClient.GetArchive(dataBytes, revision)
	if err != nil {
		return err
	}

	c.Response().Header().Set("ETag", revision)

	return c.Blob(http.StatusOK, "application/gzip", archive)
}
