package forex

import (
	"context"

	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/shopspring/decimal"
)

type Repository interface {
	GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error)
}
