package forex

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/forex/currency"
)

type ExchangeAPI interface {
	GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error)
}

type Service struct {
	api ExchangeAPI
}

func NewService(api ExchangeAPI) *Service {
	return &Service{api: api}
}

func (s *Service) GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error) {
	return s.api.GetRate(ctx, from, to)
}
