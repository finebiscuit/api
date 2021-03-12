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
	"github.com/finebiscuit/api/util"
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
	b := balance.New(currency, util.EncryptedData{
		Encryption: params.Encryption,
		Data:       params.Data,
	})

	es := make([]*entry.Entry, 0, len(params.Entries))
	for _, e := range params.Entries {
		cur := forex.Currency(params.Currency)
		es = append(es, entry.New(b.ID, cur, time.Now(), util.EncryptedData{
			Encryption: e.Encryption,
			Data:       e.Data,
		}))
	}

	if err := r.Accounting.CreateBalance(ctx, b, es); err != nil {
		return nil, err
	}

	return &model.BalancePayload{Balance: model.NewBalance(b)}, nil
}

func (r *mutationResolver) UpdateBalance(ctx context.Context, params model.UpdateBalanceInput) (*model.BalancePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveBalance(ctx context.Context, balanceID string) (*model.BalancePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddEntries(ctx context.Context, params model.AddEntriesInput) (*model.BalancePayload, error) {
	b, err := r.Accounting.FindBalance(ctx, balance.ParseID(params.BalanceID), balance.Filter{})
	if err != nil {
		return nil, err
	}

	es := make([]*entry.Entry, 0, len(params.Entries))
	for _, e := range params.Entries {
		cur := forex.Currency(e.Currency)
		es = append(es, entry.New(b.ID, cur, time.Now(), util.EncryptedData{
			Encryption: e.Encryption,
			Data:       e.Data,
		}))
	}

	if err := r.Accounting.AddEntries(ctx, es); err != nil {
		return nil, err
	}

	return &model.BalancePayload{Balance: model.NewBalance(b)}, nil
}

func (r *mutationResolver) RemoveEntry(ctx context.Context, entryID string) (*model.BalancePayload, error) {
	panic(fmt.Errorf("not implemented"))
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
