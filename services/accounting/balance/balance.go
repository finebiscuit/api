package balance

import (
	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/util"
)

//go:generate go run github.com/finebiscuit/api/scripts/genid -pkg=balance

type Balance struct {
	ID       ID
	Currency currency.Currency
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

func New(currency currency.Currency, typ Type, opt Optional) *Balance {
	return &Balance{
		Currency: currency,
		Type:     typ,
		Optional: opt,
	}
}

type ValueMap map[currency.Currency]decimal.Decimal

type WithCurrentValue struct {
	Balance

	CurrentValue ValueMap
}
