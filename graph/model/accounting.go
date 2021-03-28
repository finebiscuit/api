package model

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
)

type Balance struct {
	ID       string            `json:"id"`
	Currency currency.Currency `json:"currency"`
	Kind     string            `json:"kind"`

	DisplayName  *string `json:"displayName"`
	OfficialName *string `json:"officialName"`
	Institution  *string `json:"institution"`

	EstimatedMonthlyGrowthRate  decimal.Decimal `json:"estimatedMonthlyGrowthRate"`
	EstimatedMonthlyValueChange decimal.Decimal `json:"estimatedMonthlyValueChange"`

	ValidAt       time.Time
	CurrentValues map[currency.Currency]decimal.Decimal
}

func NewBalance(b *balance.WithCurrentValue) *Balance {
	return &Balance{
		ID:                          string(b.ID),
		Currency:                    b.Currency,
		Kind:                        b.Type.String(),
		DisplayName:                 &b.DisplayName,
		OfficialName:                &b.OfficialName,
		Institution:                 &b.Institution,
		ValidAt:                     b.ValidAt,
		CurrentValues:               b.CurrentValue,
		EstimatedMonthlyValueChange: b.EstimatedMonthlyValueChange,
		EstimatedMonthlyGrowthRate:  b.EstimatedMonthlyGrowthRate,
	}
}

func NewBalanceValue(cur currency.Currency, value decimal.Decimal, validAt time.Time) *BalanceValue {
	return &BalanceValue{
		Currency: cur,
		Value:    value,
		ValidAt:  validAt,
		Year:     validAt.Year(),
		Month:    int(validAt.Month()),
	}
}
