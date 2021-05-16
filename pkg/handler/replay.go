package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/replay"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

func (client *Client) ReplayLogWithCurrentRules(c echo.Context) error {
	decisionID := c.Param("decisionID")

	resultSet, err := client.replayClient.ReplayLog(decisionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resultSet)
}

func (client *Client) ReplayLogWithNewRules(c echo.Context) error {
	decisionID := c.Param("decisionID")
	var rules []rule.Rule

	if err := c.Bind(&rules); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tmpBundleClient := bundle.NewClient()
	tmpRuleClient := rule.NewClient()

	for _, r := range rules {
		opts := rule.Options{
			Country:    r.Country,
			City:       r.City,
			Building:   r.Building,
			Role:       r.Role,
			DeviceType: r.DeviceType,
			Action:     rule.ToAction(r.Action),
		}

		_, err := tmpRuleClient.Add(opts)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	replayOpts := replay.Options{
		RuleClient:   tmpRuleClient,
		BundleClient: tmpBundleClient,
		LogsClient:   client.logsClient,
	}

	tmpReplayClient := replay.NewClient(replayOpts)

	resultSet, err := tmpReplayClient.ReplayLog(decisionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resultSet)
}
