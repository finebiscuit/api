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
			params := model.UpdatePreferencesParams{
				DefaultCurrency: strPtr(currency.EUR.String()),
				SupportedCurrencies: []string{
					currency.PLN.String(),
					currency.RUB.String(),
				},
			}
			res, err := resolver.Mutation().UpdatePreferences(ctx, params)
			require.NoError(t, err)
			require.NotNil(t, res.Preferences)
			assert.Equal(t, params.DefaultCurrency, res.Preferences.DefaultCurrency)
			assert.Equal(t, []string{"EUR", "PLN", "RUB"}, res.Preferences.SupportedCurrencies)
		})
	})

	t.Run("Query_Preferences(AfterUpdate)", func(t *testing.T) {
		res, err := resolver.Query().Preferences(ctx)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.DefaultCurrency)
		assert.Equal(t, "EUR", *res.DefaultCurrency)
		assert.Equal(t, []string{"EUR", "PLN", "RUB"}, res.SupportedCurrencies)
	})
}
