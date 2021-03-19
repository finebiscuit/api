package entry

import (
	"time"

	"github.com/finebiscuit/api/services/forex/currency"
)

type Filter struct {
	Currencies  []currency.Currency
	ValidAfter  *time.Time
	ValidBefore *time.Time
}
