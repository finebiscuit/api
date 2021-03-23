package model

import (
	"github.com/finebiscuit/api/services/prefs"
)

type Preferences struct {
	DefaultCurrency     *string  `json:"defaultCurrency"`
	SupportedCurrencies []string `json:"supportedCurrencies"`
}

func NewPreferences(p *prefs.Preferences) *Preferences {
	var res Preferences

	if p.DefaultCurrency.IsACurrency() {
		s := p.DefaultCurrency.String()
		res.DefaultCurrency = &s
	}

	res.SupportedCurrencies = make([]string, 0, len(p.SupportedCurrencies))
	for _, c := range p.SupportedCurrencies {
		res.SupportedCurrencies = append(res.SupportedCurrencies, c.String())
	}

	return &res
}
