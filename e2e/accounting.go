package e2e

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/services/prefs"
)

func AccountingTests(t *testing.T, ctx context.Context, resolver *graph.Resolver) {
	var (
		bal *model.Balance
	)

	err := resolver.Prefs.SetPreferences(ctx, &prefs.Preferences{
		DefaultCurrency:     currency.EUR,
		SupportedCurrencies: []currency.Currency{currency.EUR, currency.RUB},
	})
	require.NoError(t, err)

	t.Run("Mutation_CreateBalance", func(t *testing.T) {
		t.Run("InvalidEmpty", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{})
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("InvalidUserUnsupportedCurrency", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{
				Currency:    currency.GBP,
				Kind:        balance.CashChecking.String(),
				Value:       decimal.RequireFromString("123.45"),
				DisplayName: strPtr("Unsupported Balance"),
			})
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("InvalidBadCurrency", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{
				Currency:    255,
				Kind:        balance.CashChecking.String(),
				Value:       decimal.RequireFromString("123.45"),
				DisplayName: strPtr("Bad Currency Balance"),
			})
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("Success", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{
				Currency:     currency.EUR,
				Kind:         balance.CashChecking.String(),
				Value:        decimal.RequireFromString("123.45"),
				DisplayName:  strPtr("Balance"),
				Institution:  strPtr("Institution"),
				OfficialName: strPtr("Institution's Balance"),
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			assert.Equal(t, currency.EUR, res.Balance.Currency)
			assert.Equal(t, "CashChecking", res.Balance.Kind)
			assert.Equal(t, "Balance", *res.Balance.DisplayName)
			assert.Equal(t, "Institution", *res.Balance.Institution)
			assert.Equal(t, "Institution's Balance", *res.Balance.OfficialName)
		})

		t.Run("SuccessNoOptional", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{
				Currency: currency.RUB,
				Kind:     balance.CashPhysical.String(),
				Value:    decimal.RequireFromString("543.21"),
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			assert.Equal(t, currency.RUB, res.Balance.Currency)
			assert.Equal(t, "CashPhysical", res.Balance.Kind)
			assert.Empty(t, res.Balance.DisplayName)
			assert.Empty(t, res.Balance.Institution)
			assert.Empty(t, res.Balance.OfficialName)
		})
	})

	t.Run("Query_Balances", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			res, err := resolver.Query().Balances(ctx)
			require.NoError(t, err)
			require.Len(t, res, 2)

			for _, b := range res {
				switch b.Currency {
				case currency.EUR:
					assert.Equal(t, "CashChecking", b.Kind)
					assert.Equal(t, "Balance", *b.DisplayName)
					assert.Equal(t, "Institution", *b.Institution)
					assert.Equal(t, "Institution's Balance", *b.OfficialName)
					bal = b
				case currency.RUB:
					assert.Equal(t, "CashPhysical", b.Kind)
					assert.Empty(t, b.DisplayName)
					assert.Empty(t, b.Institution)
					assert.Empty(t, b.OfficialName)
				}
			}
		})
	})

	t.Run("Mutation_UpdateBalanceValue", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			res, err := resolver.Mutation().UpdateBalanceValue(ctx, model.UpdateBalanceValueInput{
				BalanceID: bal.ID,
				Value:     "234.56",
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			bal = res.Balance
		})
	})

	t.Run("Balance_AllCurrentValues", func(t *testing.T) {
		res, err := resolver.Balance().AllCurrentValues(ctx, bal)
		require.NoError(t, err)
		require.Len(t, res, 1)

		assert.Equal(t, currency.EUR, res[0].Currency)
		assert.Equal(t, "234.56", res[0].Value)
	})
}

func strPtr(s string) *string {
	return &s
}
