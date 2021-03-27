package projecting

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/util"
)

type Service interface {
	ProjectBalanceValue(ctx context.Context, balanceID balance.ID, cur currency.Currency, since, until util.Period) (map[util.Period]decimal.Decimal, currency.Currency, error)
}

type service struct {
	accounting accounting.Service
}

func NewService(accountingSvc accounting.Service) Service {
	return &service{accounting: accountingSvc}
}

func (s service) ProjectBalanceValue(ctx context.Context, balanceID balance.ID, cur currency.Currency, since, until util.Period) (map[util.Period]decimal.Decimal, currency.Currency, error) {
	if since == until {
		return nil, 0, fmt.Errorf("since and until cannot be equal")
	}

	if since.Time().After(until.Time()) {
		return nil, 0, fmt.Errorf("since cannot be after until")
	}

	b, err := s.accounting.GetBalanceWithCurrentValue(ctx, balanceID)
	if err != nil {
		return nil, 0, err
	}

	if !cur.IsACurrency() {
		cur = b.Currency
	}

	val, ok := b.CurrentValue[cur]
	if !ok {
		return nil, 0, fmt.Errorf("no current value for requested currency: %s", cur)
	}

	result := make(map[util.Period]decimal.Decimal)
	for current := since.Next(); current != until.Next(); current = current.Next() {
		if !b.EstimatedMonthlyGrowthRate.IsZero() {
			val = val.Mul(b.EstimatedMonthlyGrowthRate)
		}
		val = val.Add(b.EstimatedMonthlyValueChange).Round(2)
		result[current] = val
	}
	return result, cur, nil
}
