package entry

import (
	"time"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/util"
)

//go:generate go run github.com/finebiscuit/api/scripts/genid -pkg=entry

type Entry struct {
	ID        ID
	BalanceID balance.ID
	Currency  forex.Currency
	ValidAt   time.Time
	util.EncryptedData
	util.HasTimestamps
}

func New(bID balance.ID, currency forex.Currency, validAt time.Time, ed util.EncryptedData) *Entry {
	return &Entry{
		ID:            newID(),
		BalanceID:     bID,
		Currency:      currency,
		ValidAt:       validAt,
		EncryptedData: ed,
	}
}
