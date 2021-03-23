package forex

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/forex/currency"
)

type Service interface {
	GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error)
}

type service struct {
	api Repository
}

func NewService(api Repository) Service {
	return &service{api: api}
}

func (s *service) GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error) {
	return s.api.GetRate(ctx, from, to)
}
