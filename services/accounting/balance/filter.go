package balance

import (
	"github.com/finebiscuit/api/services/forex"
)

type Filter struct {
	Currencies []forex.Currency
}
