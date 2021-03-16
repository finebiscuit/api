package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/finebiscuit/api/graph/generated"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
	"github.com/finebiscuit/api/services/forex"
	"github.com/shopspring/decimal"
)

func (r *balanceResolver) AllEntries(ctx context.Context, obj *model.Balance) ([]*model.Entry, error) {
	es, err := r.Accounting.ListBalanceEntries(ctx, balance.ParseID(obj.ID))
	if err != nil {
		return nil, err
	}

	mes := make([]*model.Entry, 0, len(es))
	for _, e := range es {
		mes = append(mes, model.NewEntry(e))
	}
	return mes, nil
}

func (r *balanceResolver) LatestEntry(ctx context.Context, obj *model.Balance, currency string) (*model.Entry, error) {
	cur := forex.Currency(currency)
	e, err := r.Accounting.FindLatestEntryByCurrency(ctx, balance.ParseID(obj.ID), cur)
	if err != nil {
		return nil, err
	}

	return model.NewEntry(e), nil
}

func (r *mutationResolver) CreateBalance(ctx context.Context, params model.CreateBalanceInput) (*model.BalancePayload, error) {
	currency := forex.Currency(params.Currency)
	typ, err := balance.TypeString(params.Kind)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %q", params.Kind)
	}

	val, err := decimal.NewFromString(params.Value)
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

	b := balance.New(currency, typ, opt)
	e := entry.New(balance.NoID, currency, val, time.Now())

	if err := r.Accounting.CreateBalance(ctx, b, []*entry.Entry{e}); err != nil {
		return nil, err
	}

	return &model.BalancePayload{Balance: model.NewBalance(b)}, nil
}

func (r *mutationResolver) RemoveBalance(ctx context.Context, balanceID string) (*model.BalancePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateBalanceInfo(ctx context.Context, params model.UpdateBalanceInfoInput) (*model.BalancePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateBalanceValue(ctx context.Context, params model.UpdateBalanceValueInput) (*model.BalancePayload, error) {
	b, err := r.Accounting.FindBalance(ctx, balance.ParseID(params.BalanceID), balance.Filter{})
	if err != nil {
		return nil, err
	}

	val, err := decimal.NewFromString(params.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %q", params.Value)
	}

	e := entry.New(b.ID, b.Currency, val, time.Now())
	if err := r.Accounting.AddEntries(ctx, []*entry.Entry{e}); err != nil {
		return nil, err
	}

	return &model.BalancePayload{Balance: model.NewBalance(b)}, nil
}

func (r *queryResolver) Balances(ctx context.Context) ([]*model.Balance, error) {
	bs, err := r.Accounting.ListBalances(ctx)
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
