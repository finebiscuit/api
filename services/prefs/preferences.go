package prefs

import (
	"github.com/finebiscuit/api/services/forex/currency"
)

type Preferences struct {
	DefaultCurrency     currency.Currency
	SupportedCurrencies []currency.Currency
}

func (p Preferences) IsCurrencySupported(cur currency.Currency) bool {
	for _, c := range p.SupportedCurrencies {
		if c == cur {
			return true
		}
	}
	return false
}

type Change struct {
	Key   Key
	Value interface{}
}
