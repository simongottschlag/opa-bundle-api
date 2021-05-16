package logs

import (
	"errors"
	"sync"

	opalogs "github.com/open-policy-agent/opa/plugins/logs"
)

var (
	NullOpaEvent         = opalogs.EventV1{}
	ErrorIDAlreadyExists = errors.New("DecisionID already exists")
	ErrorIDNotFound      = errors.New("DecisionID not found")
)

type DecisionID = string

type Client struct {
	sync.RWMutex
	logs map[DecisionID]opalogs.EventV1
}

func NewClient() *Client {
	return &Client{
		logs: make(map[DecisionID]opalogs.EventV1),
	}
}

func (client *Client) Create(log opalogs.EventV1) error {
	client.Lock()
	defer client.Unlock()

	_, found := client.logs[log.DecisionID]
	if found {
		return ErrorIDAlreadyExists
	}

	client.logs[log.DecisionID] = log

	return nil
}

func (client *Client) CreateMultiple(logs []opalogs.EventV1) error {
	client.Lock()
	defer client.Unlock()

	for _, log := range logs {
		err := client.createWithoutLock(log)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) createWithoutLock(log opalogs.EventV1) error {
	_, found := client.logs[log.DecisionID]
	if found {
		return ErrorIDAlreadyExists
	}

	client.logs[log.DecisionID] = log

	return nil
}

func (client *Client) Read(id DecisionID) (opalogs.EventV1, error) {
	client.RLock()
	defer client.RUnlock()

	log, found := client.logs[id]
	if !found {
		return NullOpaEvent, ErrorIDNotFound
	}

	return log, nil
}

func (client *Client) ReadAll() []opalogs.EventV1 {
	client.RLock()
	defer client.RUnlock()

	var logs []opalogs.EventV1
	for _, log := range client.logs {
		logs = append(logs, log)
	}

	return logs
}
