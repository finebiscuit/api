package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/finebiscuit/api/graph/generated"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	forexcurrency "github.com/finebiscuit/api/services/forex/currency"
)

func (r *balanceResolver) AllCurrentValues(ctx context.Context, obj *model.Balance) ([]*model.BalanceValue, error) {
	mvals := make([]*model.BalanceValue, 0, len(obj.CurrentValues))
	for k, v := range obj.CurrentValues {
		mvals = append(mvals, model.NewBalanceValue(k, v, obj.ValidAt))
	}
	return mvals, nil
}

func (r *balanceResolver) CurrentValue(ctx context.Context, obj *model.Balance, currency forexcurrency.Currency) (*model.BalanceValue, error) {
	// TODO: in case there's no value for this currency, use the forex service to dynamically calculate it.
	v, ok := obj.CurrentValues[currency]
	if !ok {
		return nil, fmt.Errorf("no value found for currency %q", currency)
	}

	return model.NewBalanceValue(currency, v, obj.ValidAt), nil
}

func (r *balanceResolver) HistoricalValues(ctx context.Context, obj *model.Balance, currency *forexcurrency.Currency) ([]*model.BalanceValue, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateBalance(ctx context.Context, params model.CreateBalanceInput) (*model.BalancePayload, error) {
	typ, err := balance.TypeString(params.Kind)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %q", params.Kind)
	}

	opt := balance.Optional{}
	if params.DisplayName != nil {
		opt.DisplayName = *params.DisplayName
	}
	if params.OfficialName != nil {
		opt.OfficialName = *params.OfficialName
	}
	if params.Institution != nil {
		opt.Institution = *params.Institution
	}
	if params.EstimatedMonthlyGrowthRate != nil {
		opt.EstimatedMonthlyGrowthRate = *params.EstimatedMonthlyGrowthRate
	}
	if params.EstimatedMonthlyValueChange != nil {
		opt.EstimatedMonthlyValueChange = *params.EstimatedMonthlyValueChange
	}

	b := balance.New(params.Currency, typ, opt)

	values, err := r.Accounting.CreateBalance(ctx, b, params.Value)
	if err != nil {
		return nil, err
	}

	payload := &model.BalancePayload{
		Balance: model.NewBalance(&balance.WithCurrentValue{
			Balance:      *b,
			CurrentValue: values,
		}),
	}
	return payload, nil
}

func (r *mutationResolver) RemoveBalance(ctx context.Context, balanceID string) (*model.BalancePayload, error) {
	b, err := r.Accounting.GetBalance(ctx, balance.ParseID(balanceID))
	if err != nil {
		return nil, err
	}

	if err := r.Accounting.DeleteBalance(ctx, b.ID); err != nil {
		return nil, err
	}
	return &model.BalancePayload{}, nil
}

func (r *mutationResolver) UpdateBalanceInfo(ctx context.Context, params model.UpdateBalanceInfoInput) (*model.BalancePayload, error) {
	b, err := r.Accounting.GetBalanceWithCurrentValue(ctx, balance.ParseID(params.BalanceID))
	if err != nil {
		return nil, err
	}

	if params.DisplayName != nil {
		b.DisplayName = *params.DisplayName
	}
	if params.OfficialName != nil {
		b.OfficialName = *params.OfficialName
	}
	if params.Institution != nil {
		b.Institution = *params.Institution
	}
	if params.EstimatedMonthlyGrowthRate != nil {
		b.EstimatedMonthlyGrowthRate = *params.EstimatedMonthlyGrowthRate
	}
	if params.EstimatedMonthlyValueChange != nil {
		b.EstimatedMonthlyValueChange = *params.EstimatedMonthlyValueChange
	}

	if err := r.Accounting.UpdateBalanceInfo(ctx, &b.Balance); err != nil {
		return nil, err
	}

	return &model.BalancePayload{Balance: model.NewBalance(b)}, nil
}

func (r *mutationResolver) UpdateBalanceValue(ctx context.Context, params model.UpdateBalanceValueInput) (*model.BalancePayload, error) {
	b, err := r.Accounting.GetBalance(ctx, balance.ParseID(params.BalanceID))
	if err != nil {
		return nil, err
	}

	values, err := r.Accounting.UpdateBalanceValue(ctx, b.ID, params.Value)
	if err != nil {
		return nil, err
	}

	payload := &model.BalancePayload{
		Balance: model.NewBalance(&balance.WithCurrentValue{
			Balance:      *b,
			CurrentValue: values,
		}),
	}
	return payload, nil
}

func (r *queryResolver) Balances(ctx context.Context) ([]*model.Balance, error) {
	bs, err := r.Accounting.ListBalancesWithCurrentValue(ctx)
	if err != nil {
		return nil, err
	}

	mbs := make([]*model.Balance, 0, len(bs))
	for _, b := range bs {
		mbs = append(mbs, model.NewBalance(b))
	}

	return mbs, nil
}

func (r *queryResolver) Balance(ctx context.Context, id string) (*model.Balance, error) {
	b, err := r.Accounting.GetBalanceWithCurrentValue(ctx, balance.ParseID(id))
	if err != nil {
		return nil, err
	}
	return model.NewBalance(b), nil
}

// Balance returns generated.BalanceResolver implementation.
func (r *Resolver) Balance() generated.BalanceResolver { return &balanceResolver{r} }

type balanceResolver struct{ *Resolver }
