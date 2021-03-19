package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/finebiscuit/api/graph/generated"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/shopspring/decimal"
)

func (r *balanceResolver) AllCurrentValues(ctx context.Context, obj *model.Balance) ([]*model.BalanceValue, error) {
	mvals := make([]*model.BalanceValue, 0, len(obj.CurrentValues))
	for k, v := range obj.CurrentValues {
		mvals = append(mvals, &model.BalanceValue{Currency: k.String(), Value: v.String()})
	}
	return mvals, nil
}

func (r *balanceResolver) CurrentValue(ctx context.Context, obj *model.Balance, cur string) (*model.BalanceValue, error) {
	c := currency.New(cur)

	// TODO: in case there's no value for this currency, use the forex service to dynamically calculate it.
	v, ok := obj.CurrentValues[c]
	if !ok {
		return nil, fmt.Errorf("no value found for currency %q", c)
	}

	mval := &model.BalanceValue{Currency: c.String(), Value: v.String()}
	return mval, nil
}

func (r *mutationResolver) CreateBalance(ctx context.Context, params model.CreateBalanceInput) (*model.BalancePayload, error) {
	cur := currency.New(params.Currency)
	typ, err := balance.TypeString(params.Kind)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %q", params.Kind)
	}

	decVal, err := decimal.NewFromString(params.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %q", params.Value)
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

	b := balance.New(cur, typ, opt)
	values := map[currency.Currency]decimal.Decimal{
		b.Currency: decVal,
	}

	if err := r.Accounting.CreateBalance(ctx, b, values); err != nil {
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

	decVal, err := decimal.NewFromString(params.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %q", params.Value)
	}

	values := map[currency.Currency]decimal.Decimal{
		b.Currency: decVal,
	}

	if err := r.Accounting.UpdateBalanceValue(ctx, b.ID, values); err != nil {
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

func (r *queryResolver) Preferences(ctx context.Context) (*model.Preferences, error) {
	panic(fmt.Errorf("not implemented"))
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

// Balance returns generated.BalanceResolver implementation.
func (r *Resolver) Balance() generated.BalanceResolver { return &balanceResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type balanceResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
