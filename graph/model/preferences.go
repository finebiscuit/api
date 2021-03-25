package model

import (
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/services/prefs"
)

type Preferences struct {
	DefaultCurrency     *currency.Currency  `json:"defaultCurrency"`
	SupportedCurrencies []currency.Currency `json:"supportedCurrencies"`
}

func NewPreferences(p *prefs.Preferences) *Preferences {
	var res Preferences

	if p.DefaultCurrency.IsACurrency() {
		res.DefaultCurrency = &p.DefaultCurrency
	}

	if len(p.SupportedCurrencies) > 0 {
		res.SupportedCurrencies = p.SupportedCurrencies
	}

	return &res
}
