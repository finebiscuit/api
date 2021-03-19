package model

import (
	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
)

type Balance struct {
	ID           string  `json:"id"`
	Currency     string  `json:"currency"`
	Kind         string  `json:"kind"`
	DisplayName  *string `json:"displayName"`
	OfficialName *string `json:"officialName"`
	Institution  *string `json:"institution"`

	CurrentValues map[currency.Currency]decimal.Decimal
}

func NewBalance(b *balance.WithCurrentValue) *Balance {
	return &Balance{
		ID:            string(b.ID),
		Currency:      b.Currency.String(),
		Kind:          b.Type.String(),
		DisplayName:   &b.DisplayName,
		OfficialName:  &b.OfficialName,
		Institution:   &b.Institution,
		CurrentValues: b.CurrentValue,
	}
}
