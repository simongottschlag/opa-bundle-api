package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

func (client *Client) ReadRules(c echo.Context) error {
	rules, err := client.repository.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Bind(rules); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, rules)
}

func (client *Client) CreateRule(c echo.Context) error {
	r := rule.Rule{}

	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	opts := rule.Options{
		Country:    r.Country,
		City:       r.City,
		Building:   r.Building,
		Role:       r.Role,
		DeviceType: r.DeviceType,
		Action:     rule.ToAction(r.Action),
	}

	id, err := client.repository.Add(opts)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rule, err := client.repository.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, rule)
}

func (client *Client) ReadRule(c echo.Context) error {
	id, err := rule.StringToID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rule, err := client.repository.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, rule)
}

func (client *Client) UpdateRule(c echo.Context) error {
	id, err := rule.StringToID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	r := rule.Rule{}

	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	opts := rule.Options{
		Country:    r.Country,
		City:       r.City,
		Building:   r.Building,
		Role:       r.Role,
		DeviceType: r.DeviceType,
		Action:     rule.ToAction(r.Action),
	}

	err = client.repository.Set(id, opts)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rule, err := client.repository.Get(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, rule)
}

func (client *Client) DeleteRule(c echo.Context) error {
	id, err := rule.StringToID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = client.repository.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
