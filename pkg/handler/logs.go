package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	opalogs "github.com/open-policy-agent/opa/plugins/logs"
)

func (client *Client) CreateLogs(c echo.Context) error {
	var logs []opalogs.EventV1

	if err := c.Bind(&logs); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := client.logsClient.CreateMultiple(logs)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (client *Client) ReadLogs(c echo.Context) error {
	logs := client.logsClient.ReadAll()

	if err := c.Bind(logs); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, logs)
}

func (client *Client) ReadLog(c echo.Context) error {
	decisionID := c.Param("decisionID")

	log, err := client.logsClient.Read(decisionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, log)
}
