package entry

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/util"
)

//go:generate go run github.com/finebiscuit/api/scripts/genid -pkg=entry

type Entry struct {
	ID        ID
	BalanceID balance.ID
	Currency  currency.Currency
	Value     decimal.Decimal
	ValidAt   time.Time

	util.HasTimestamps
}

func New(bID balance.ID, currency currency.Currency, value decimal.Decimal, validAt time.Time) *Entry {
	return &Entry{
		BalanceID: bID,
		Currency:  currency,
		Value:     value,
		ValidAt:   validAt,
	}
}
