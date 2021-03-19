package balance

import (
	"github.com/finebiscuit/api/services/forex/currency"
)

type Filter struct {
	Currencies []currency.Currency
}
