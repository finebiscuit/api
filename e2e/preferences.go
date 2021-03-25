package e2e

import (
	"context"
	"testing"

	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func PreferencesTests(t *testing.T, ctx context.Context, resolver *graph.Resolver) {
	t.Run("Query_Preferences(FirstTime)", func(t *testing.T) {
		res, err := resolver.Query().Preferences(ctx)
		require.NoError(t, err)
		assert.Nil(t, res.DefaultCurrency)
		assert.Empty(t, res.SupportedCurrencies)
	})

	t.Run("Mutation_UpdatePreferences", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			cur := currency.EUR
			params := model.UpdatePreferencesParams{
				DefaultCurrency: &cur,
				SupportedCurrencies: []currency.Currency{
					currency.PLN,
					currency.RUB,
				},
			}
			res, err := resolver.Mutation().UpdatePreferences(ctx, params)
			require.NoError(t, err)
			require.NotNil(t, res.Preferences)
			assert.Equal(t, params.DefaultCurrency, res.Preferences.DefaultCurrency)
			assert.Len(t, res.Preferences.SupportedCurrencies, 3)
			assert.Contains(t, res.Preferences.SupportedCurrencies, currency.EUR)
			assert.Contains(t, res.Preferences.SupportedCurrencies, currency.PLN)
			assert.Contains(t, res.Preferences.SupportedCurrencies, currency.RUB)
		})
	})

	t.Run("Query_Preferences(AfterUpdate)", func(t *testing.T) {
		res, err := resolver.Query().Preferences(ctx)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.DefaultCurrency)
		assert.Len(t, res.SupportedCurrencies, 3)
		assert.Contains(t, res.SupportedCurrencies, currency.EUR)
		assert.Contains(t, res.SupportedCurrencies, currency.PLN)
		assert.Contains(t, res.SupportedCurrencies, currency.RUB)
	})
}
