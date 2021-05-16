package replay

import (
	"context"

	"github.com/open-policy-agent/opa/rego"
	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/logs"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
	"github.com/xenitab/opa-bundle-api/pkg/util"
)

var (
	NullOpaResultSet = rego.ResultSet{}
)

type Options struct {
	RuleClient   *rule.Client
	BundleClient *bundle.Client
	LogsClient   *logs.Client
}

type Client struct {
	bundleClient *bundle.Client
	logsClient   *logs.Client
	ruleClient   *rule.Client
}

func NewClient(opts Options) *Client {
	return &Client{
		ruleClient:   opts.RuleClient,
		bundleClient: opts.BundleClient,
		logsClient:   opts.LogsClient,
	}
}

func (client *Client) ReplayLog(decisionID string) (rego.ResultSet, error) {
	log, err := client.logsClient.Read(decisionID)
	if err != nil {
		return NullOpaResultSet, err
	}

	input := *log.Input

	data, err := client.ruleClient.GetAllJSON()
	if err != nil {
		return NullOpaResultSet, err
	}

	dataBytes := []byte(data)
	revision, err := util.BytesToHash(dataBytes)
	if err != nil {
		return NullOpaResultSet, err
	}

	bundle, err := client.bundleClient.Get(dataBytes, revision)
	if err != nil {
		return NullOpaResultSet, err
	}

	ctx := context.Background()

	rego := rego.New(
		rego.ParsedBundle("bundle", &bundle),
		rego.Input(input),
		rego.Query(`data.rule.allow`),
	)

	resultSet, err := rego.Eval(ctx)
	if err != nil {
		return NullOpaResultSet, err
	}

	return resultSet, nil
}
