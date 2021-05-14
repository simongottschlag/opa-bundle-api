package bundle

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

type Bundle struct {
	Manifest manifest
	Rules    *rule.Rules
	Policies Policies
}

type manifest struct {
	Revision string   `json:"revision"`
	Roots    []string `json:"roots,omitempty"`
}

func (m *manifest) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *manifest) UnmarshalJSON(b []byte) error {
	var res manifest
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}

	*m = res

	return nil
}

type Policy struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Policies []Policy

func (p *Policies) MarshalJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Policies) UnmarshalJSON(b []byte) error {
	var res Policies
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}

	*p = res

	return nil
}

func GenerateBundle(r *rule.Rules, p Policies) (Bundle, error) {
	pol, err := p.MarshalJSON()
	if err != nil {
		return Bundle{}, err
	}

	polHash, err := getHash(pol)
	if err != nil {
		return Bundle{}, err
	}

	rulesJson, err := r.GetAllJSON()
	if err != nil {
		return Bundle{}, err
	}

	rulesHash, err := getHash([]byte(rulesJson))
	if err != nil {
		return Bundle{}, err
	}

	revision, err := getHash([]byte(fmt.Sprintf("%s-%s", polHash, rulesHash)))
	if err != nil {
		return Bundle{}, err
	}

	m := manifest{
		Revision: revision,
	}

	b := Bundle{
		Manifest: m,
		Rules:    r,
		Policies: p,
	}

	return b, nil
}

func getHash(b []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
