package rule

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"sync"
)

var (
	WildcardString           = "ANY"
	NullRule                 = Rule{}
	NullRuleString           = ""
	NullID                   = 0
	NullAction               = ""
	ErrorIdAlreadyExists     = errors.New("ID already exists")
	ErrorIdNotFound          = errors.New("ID not found")
	ErrorNotAbleToGenerateId = errors.New("Not able to generate ID")
	ErrorNotAbleToParseId    = errors.New("Not able to parse ID")
	ErrorUnableToMarshalJSON = errors.New("Unable to marshal JSON")
	ErrorRuleNotValid        = errors.New("Rule not valid")
)

type ID = int

type Action int

const (
	ActionUndefined Action = iota
	ActionAllow
	ActionDeny
)

type Options struct {
	Country    string
	City       string
	Building   string
	Role       string
	DeviceType string
	Action     Action
}

type Rule struct {
	ID         ID     `json:"id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Building   string `json:"building"`
	Role       string `json:"role"`
	DeviceType string `json:"device_type"`
	Action     string `json:"action"`
}

func (r *Rule) Valid() bool {
	if isEmpty(r.Country, r.City, r.Building, r.Role, r.Action, r.DeviceType) {
		return false
	}

	return true

}

type Rules struct {
	sync.RWMutex
	Index      int
	Repository map[ID]Rule
}

func NewRepository() *Rules {
	return &Rules{
		Repository: make(map[ID]Rule),
	}
}

func StringToID(id string) (ID, error) {
	return strconv.Atoi(id)
}

func (r *Rules) Add(opts Options) (ID, error) {
	r.Lock()
	defer r.Unlock()

	id := r.Index
	id++

	_, found := r.Repository[id]
	if found {
		return NullID, ErrorIdAlreadyExists
	}

	rule := Rule{
		ID:         id,
		Country:    opts.Country,
		City:       opts.City,
		Building:   opts.Building,
		Role:       opts.Role,
		DeviceType: opts.DeviceType,
		Action:     FromAction(opts.Action),
	}

	if !rule.Valid() {
		return NullID, ErrorRuleNotValid
	}

	r.Repository[id] = rule
	r.Index++

	return id, nil
}

func (r *Rules) Get(id ID) (Rule, error) {
	r.RLock()
	defer r.RUnlock()

	rule, found := r.Repository[id]
	if !found {
		return NullRule, ErrorIdNotFound
	}

	return rule, nil
}

func (r *Rules) GetJSON(id ID) (string, error) {
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

	var ids []int
	for k := range r.Repository {
		ids = append(ids, k)
	}

	sort.Ints(ids)

	var rules []Rule
	for _, v := range ids {
		rule := r.Repository[v]
		rules = append(rules, rule)
	}

	return rules, nil
}

func (r *Rules) GetAllJSON() (string, error) {
	r.RLock()
	defer r.RUnlock()

	var obj struct {
		Rules []Rule `json:"rules"`
	}

	var ids []int
	for k := range r.Repository {
		ids = append(ids, k)
	}

	sort.Ints(ids)

	for _, v := range ids {
		rule := r.Repository[v]
		obj.Rules = append(obj.Rules, rule)
	}

	res, err := json.Marshal(&obj)
	if err != nil {
		return NullRuleString, ErrorUnableToMarshalJSON
	}

	return string(res), nil
}

func (r *Rules) Set(id ID, opts Options) error {
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

	if !isEmpty(opts.DeviceType) {
		rule.DeviceType = opts.DeviceType
	}

	if FromAction(opts.Action) != "undefined" {
		rule.Action = FromAction(opts.Action)
	}

	r.Repository[id] = rule

	return nil
}

func (r *Rules) Delete(id ID) error {
	r.Lock()
	defer r.Unlock()

	_, found := r.Repository[id]
	if !found {
		return ErrorIdNotFound
	}

	delete(r.Repository, id)

	return nil
}

func FromAction(action Action) string {
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

func ToAction(action string) Action {
	switch action {
	case "allow":
		return ActionAllow
	case "deny":
		return ActionDeny
	case "undefined":
		return ActionUndefined
	default:
		return ActionUndefined
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
