package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (client *Client) ReplayLog(c echo.Context) error {
	decisionID := c.Param("decisionID")

	resultSet, err := client.replayClient.ReplayLog(decisionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resultSet)
}
