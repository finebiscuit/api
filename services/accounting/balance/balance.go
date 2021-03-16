package balance

import (
	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/util"
)

//go:generate go run github.com/finebiscuit/api/scripts/genid -pkg=balance

type Balance struct {
	ID       ID
	Currency forex.Currency
	Type     Type

	Optional

	util.HasTimestamps
}

type Optional struct {
	DisplayName  string
	OfficialName string
	Institution  string

	EstimatedMonthlyGrowthRate  decimal.Decimal
	EstimatedMonthlyValueChange decimal.Decimal
}

func New(currency forex.Currency, typ Type, opt Optional) *Balance {
	return &Balance{
		Currency: currency,
		Type:     typ,
		Optional: opt,
	}
}
