package balance

import (
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/util"
)

//go:generate go run github.com/finebiscuit/api/scripts/genid -pkg=balance

type Balance struct {
	ID       ID
	Currency forex.Currency
	util.EncryptedData
	util.HasTimestamps
}

func New(currency forex.Currency, ed util.EncryptedData) *Balance {
	return &Balance{
		ID:            newID(),
		Currency:      currency,
		EncryptedData: ed,
	}
}
