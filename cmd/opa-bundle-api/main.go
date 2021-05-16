package main

import (
	"fmt"
	"net"
	"os"

	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/config"
	"github.com/xenitab/opa-bundle-api/pkg/handler"
	"github.com/xenitab/opa-bundle-api/pkg/logs"
	"github.com/xenitab/opa-bundle-api/pkg/rule"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	// Version is set at build time to print the released version using --version
	Version = "v0.0.0-dev"
	// Revision is set at build time to print the release git commit sha using --version
	Revision = ""
	// Created is set at build time to print the timestamp for when it was built using --version
	Created = ""
)

func main() {
	cfg, err := newConfigClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate config: %q\n", err)
		os.Exit(1)
	}

	err = start(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func start(cfg config.Client) error {
	rules := rule.NewRepository()

	err := seedRepository(rules)
	if err != nil {
		return err
	}

	bundleClient := bundle.NewClient()
	logsClient := logs.NewClient()
	handlerClient := newHandlerClient(rules, bundleClient, logsClient)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())

	e.GET("/", handlerClient.Default)

	eRules := e.Group("/rules")
	eRules.GET("", handlerClient.ReadRules)
	eRules.POST("", handlerClient.CreateRule)
	eRules.GET("/:id", handlerClient.ReadRule)
	eRules.PUT("/:id", handlerClient.UpdateRule)
	eRules.DELETE("/:id", handlerClient.DeleteRule)

	eLogs := e.Group("/logs")
	eLogs.POST("", handlerClient.CreateLogs, middleware.Decompress())
	eLogs.GET("", handlerClient.ReadLogs)
	eLogs.GET("/:decisionID", handlerClient.ReadLog)

	eBundle := e.Group("/bundle")
	eBundle.GET("/bundle.tar.gz", handlerClient.GetBundle)

	address := net.JoinHostPort(cfg.Address, fmt.Sprintf("%d", cfg.Port))
	e.Logger.Fatal(e.Start(address))

	return nil
}

func newConfigClient() (config.Client, error) {
	opts := config.Options{
		Version:  Version,
		Revision: Revision,
		Created:  Created,
	}

	return config.NewClient(opts)
}

func newHandlerClient(repository *rule.Rules, bundleClient *bundle.Client, logsClient *logs.Client) *handler.Client {
	opts := handler.Options{
		Repository:   repository,
		BundleClient: bundleClient,
		LogsClient:   logsClient,
	}

	return handler.NewClient(opts)
}

func seedRepository(repostiory *rule.Rules) error {
	rules := []rule.Options{
		{
			// super_admin should have access to everything
			Country:    rule.WildcardString,
			City:       rule.WildcardString,
			Building:   rule.WildcardString,
			Role:       "super_admin",
			DeviceType: rule.WildcardString,
			Action:     rule.ActionAllow,
		},
		{
			// sweden_admin should have access to everything in Sweden
			Country:    "Sweden",
			City:       rule.WildcardString,
			Building:   rule.WildcardString,
			Role:       "sweden_admin",
			DeviceType: rule.WildcardString,
			Action:     rule.ActionAllow,
		},
		{
			// norway_admin should have access to everything in Norway
			Country:    "Norway",
			City:       rule.WildcardString,
			Building:   rule.WildcardString,
			Role:       "norway_admin",
			DeviceType: rule.WildcardString,
			Action:     rule.ActionAllow,
		},
		{
			// printer_admin should have access to all Printers
			Country:    rule.WildcardString,
			City:       rule.WildcardString,
			Building:   rule.WildcardString,
			Role:       "printer_admin",
			DeviceType: "Printer",
			Action:     rule.ActionAllow,
		},
		{
			// user should have access to all Printers in Branch (Alingsås, Sweden)
			Country:    "Sweden",
			City:       "Alingsås",
			Building:   "Branch",
			Role:       "user",
			DeviceType: "Printer",
			Action:     rule.ActionAllow,
		},
		{
			// sweden_manager should have access to all Printers in Sweden
			Country:    "Sweden",
			City:       rule.WildcardString,
			Building:   rule.WildcardString,
			Role:       "sweden_manager",
			DeviceType: "Printer",
			Action:     rule.ActionAllow,
		},
		{
			// janitor should have access to Alarm in HQ (Gothenburg, Sweden)
			Country:    "Sweden",
			City:       "Gothenburg",
			Building:   "HQ",
			Role:       "janitor",
			DeviceType: "Alarm",
			Action:     rule.ActionAllow,
		},
		{
			// janitor should have access to all Alarms in Alingsås (Sweden)
			Country:    "Sweden",
			City:       "Alingsås",
			Building:   rule.WildcardString,
			Role:       "janitor",
			DeviceType: "Alarm",
			Action:     rule.ActionAllow,
		},
		{
			// guests should be denied everything
			Country:    rule.WildcardString,
			City:       rule.WildcardString,
			Building:   rule.WildcardString,
			Role:       "guest",
			DeviceType: rule.WildcardString,
			Action:     rule.ActionDeny,
		},
	}

	for _, opt := range rules {
		_, err := repostiory.Add(opt)
		if err != nil {
			return err
		}
	}

	return nil
}
