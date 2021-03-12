package entry

import (
	"time"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex"
)

type Filter struct {
	BalanceIDs  []balance.ID
	Currencies  []forex.Currency
	ValidAfter  *time.Time
	ValidBefore *time.Time
}
