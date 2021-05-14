package rule

import (
	"encoding/json"
	"errors"
	"sync"

	uuid "github.com/google/uuid"
)

var (
	NullRule                 = Rule{}
	NullRuleString           = ""
	NullID                   = uuid.Nil
	NullAction               = ""
	ErrorIdNil               = errors.New("ID can't be nil")
	ErrorIdAlreadyExists     = errors.New("ID already exists")
	ErrorIdNotFound          = errors.New("ID not found")
	ErrorUnableToMarshalJSON = errors.New("Unable to marshal JSON")
	ErrorRuleNotValid        = errors.New("Rule not valid")
)

type ID = uuid.UUID

type Action int

const (
	ActionUndefined Action = iota
	ActionAllow
	ActionDeny
)

func NewID() ID {
	id := uuid.New()
	return id
}

type Options struct {
	Country  string
	City     string
	Building string
	Role     string
	Action   Action
}

type Rule struct {
	ID       ID     `json:"id"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Building string `json:"building"`
	Role     string `json:"role"`
	Action   string `json:"action"`
}

func (r *Rule) Valid() bool {
	if isEmpty(r.Country, r.City, r.Building, r.Role, r.Action) {
		return false
	}

	if r.ID == NullID {
		return false
	}

	return true

}

type Rules struct {
	sync.RWMutex
	Repository map[ID]Rule
}

func NewRepository() Rules {
	return Rules{
		Repository: make(map[ID]Rule),
	}
}

func (r *Rules) Add(opts Options) (ID, error) {
	id := NewID()

	r.Lock()
	defer r.Unlock()

	_, found := r.Repository[id]
	if found {
		return NullID, ErrorIdAlreadyExists
	}

	rule := Rule{
		ID:       id,
		Country:  opts.Country,
		City:     opts.City,
		Building: opts.Building,
		Role:     opts.Role,
		Action:   newAction(opts.Action),
	}

	if !rule.Valid() {
		return NullID, ErrorRuleNotValid
	}

	r.Repository[id] = rule

	return id, nil
}

func (r *Rules) Get(id uuid.UUID) (Rule, error) {
	if id == uuid.Nil {
		return NullRule, ErrorIdNil
	}

	r.RLock()
	defer r.RUnlock()

	rule, found := r.Repository[id]
	if !found {
		return NullRule, ErrorIdNotFound
	}

	return rule, nil
}

func (r *Rules) GetJSON(id uuid.UUID) (string, error) {
	if id == uuid.Nil {
		return NullRuleString, ErrorIdNil
	}

	r.RLock()
	defer r.RUnlock()

	rule, found := r.Repository[id]
	if !found {
		return NullRuleString, ErrorIdNotFound
	}

	res, err := json.Marshal(&rule)
	if err != nil {
		return NullRuleString, ErrorUnableToMarshalJSON
	}

	return string(res), nil
}

func (r *Rules) GetAll() ([]Rule, error) {
	r.RLock()
	defer r.RUnlock()

	var rules []Rule
	for _, rule := range r.Repository {
		rules = append(rules, rule)
	}

	return rules, nil
}

func (r *Rules) GetAllJSON() (string, error) {
	r.RLock()
	defer r.RUnlock()

	var rules []Rule
	for _, rule := range r.Repository {
		rules = append(rules, rule)
	}

	res, err := json.Marshal(&rules)
	if err != nil {
		return NullRuleString, ErrorUnableToMarshalJSON
	}

	return string(res), nil
}

func (r *Rules) Set(id uuid.UUID, opts Options) error {
	if id == uuid.Nil {
		return ErrorIdNil
	}

	r.Lock()
	defer r.Unlock()

	rule := r.Repository[id]

	if !isEmpty(opts.Country) {
		rule.Country = opts.Country
	}

	if !isEmpty(opts.City) {
		rule.City = opts.City
	}

	if !isEmpty(opts.Building) {
		rule.Building = opts.Building
	}

	if !isEmpty(opts.Role) {
		rule.Role = opts.Role
	}

	if newAction(opts.Action) != "undefined" {
		rule.Action = newAction(opts.Action)
	}

	r.Repository[id] = rule

	return nil
}

func (r *Rules) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrorIdNil
	}

	r.Lock()
	defer r.Unlock()

	delete(r.Repository, id)

	return nil
}

func newAction(action Action) string {
	switch action {
	case ActionAllow:
		return "allow"
	case ActionDeny:
		return "deny"
	case ActionUndefined:
		return "undefined"
	default:
		return "undefined"
	}
}

func isEmpty(input ...string) bool {
	for _, s := range input {
		if s == "" {
			return true
		}
	}

	return false
}
