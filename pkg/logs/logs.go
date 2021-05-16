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

type Logs struct {
	sync.RWMutex
	logs map[DecisionID]opalogs.EventV1
}

func NewClient() *Logs {
	return &Logs{
		logs: make(map[DecisionID]opalogs.EventV1),
	}
}

func (l *Logs) Create(log opalogs.EventV1) error {
	l.Lock()
	defer l.Unlock()

	_, found := l.logs[log.DecisionID]
	if found {
		return ErrorIDAlreadyExists
	}

	l.logs[log.DecisionID] = log

	return nil
}

func (l *Logs) CreateMultiple(logs []opalogs.EventV1) error {
	l.Lock()
	defer l.Unlock()

	for _, log := range logs {
		err := l.createWithoutLock(log)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Logs) createWithoutLock(log opalogs.EventV1) error {
	_, found := l.logs[log.DecisionID]
	if found {
		return ErrorIDAlreadyExists
	}

	l.logs[log.DecisionID] = log

	return nil
}

func (l *Logs) Read(id DecisionID) (opalogs.EventV1, error) {
	l.RLock()
	defer l.RUnlock()

	log, found := l.logs[id]
	if !found {
		return NullOpaEvent, ErrorIDNotFound
	}

	return log, nil
}

func (l *Logs) ReadAll() []opalogs.EventV1 {
	l.RLock()
	defer l.RUnlock()

	var logs []opalogs.EventV1
	for _, log := range l.logs {
		logs = append(logs, log)
	}

	return logs
}
