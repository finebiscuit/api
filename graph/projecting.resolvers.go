package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/util"
)

func (r *balanceResolver) ProjectedValues(ctx context.Context, obj *model.Balance, currency currency.Currency, forMonths int) ([]*model.BalanceValue, error) {
	balanceID := balance.ParseID(obj.ID)
	since := util.NewPeriodFromTime(time.Now())
	until := util.NewPeriodFromTime(time.Now().AddDate(0, forMonths, 0))

	valMap, err := r.Projecting.ProjectBalanceValue(ctx, balanceID, currency, since, until)
	if err != nil {
		return nil, err
	}

	values := make([]*model.BalanceValue, 0, len(valMap))
	for current := since.Next(); current != until.Next(); current = current.Next() {
		values = append(values, model.NewBalanceValue(currency, valMap[current], current.Time()))
	}
	return values, nil
}
