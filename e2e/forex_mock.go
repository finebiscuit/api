package e2e

import (
	"context"

	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/shopspring/decimal"
)

type forexMock struct{}

func (forexMock) GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error) {
	return decimal.NewFromInt(1), nil
}
