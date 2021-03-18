package entry

import (
	"time"

	"github.com/finebiscuit/api/services/forex"
)

type Filter struct {
	Currencies  []forex.Currency
	ValidAfter  *time.Time
	ValidBefore *time.Time
}
