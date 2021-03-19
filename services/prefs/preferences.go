package prefs

import (
	"github.com/finebiscuit/api/services/forex/currency"
)

type Preferences struct {
	DefaultCurrency     currency.Currency
	SupportedCurrencies []currency.Currency
}
