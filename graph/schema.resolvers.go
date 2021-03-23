package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/finebiscuit/api/graph/generated"
	"github.com/finebiscuit/api/graph/model"
	forexcurrency "github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/services/prefs"
)

func (r *mutationResolver) UpdatePreferences(ctx context.Context, params model.UpdatePreferencesParams) (*model.PreferencesPayload, error) {
	var p prefs.Preferences

	if params.DefaultCurrency != nil {
		cur, err := forexcurrency.CurrencyString(*params.DefaultCurrency)
		if err != nil {
			return nil, err
		}
		p.DefaultCurrency = cur
	}

	if params.SupportedCurrencies != nil {
		p.SupportedCurrencies = make([]forexcurrency.Currency, 0, len(params.SupportedCurrencies))
		for _, s := range params.SupportedCurrencies {
			c, err := forexcurrency.CurrencyString(s)
			if err != nil {
				return nil, err
			}
			p.SupportedCurrencies = append(p.SupportedCurrencies, c)
		}
	}

	if err := r.Prefs.SetPreferences(ctx, &p); err != nil {
		return nil, err
	}

	payload := &model.PreferencesPayload{Preferences: model.NewPreferences(&p)}
	return payload, nil
}

func (r *queryResolver) Version(ctx context.Context) (*model.Version, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Preferences(ctx context.Context) (*model.Preferences, error) {
	p, err := r.Prefs.GetPreferences(ctx)
	if err != nil {
		return nil, err
	}
	return model.NewPreferences(p), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
