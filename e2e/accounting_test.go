package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/sqldb"
)

func accountingTests(t *testing.T, ctx context.Context, cfg *config.Config) {
	resolver, err := graph.NewResolver(cfg, sqldb.NewBackend())
	require.NoError(t, err)

	var (
		balanceID string
	)

	t.Run("Mutation_CreateBalance", func(t *testing.T) {
		t.Run("InvalidEmpty", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{})
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("Success", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{
				Currency:     "EUR",
				Kind:         balance.CashChecking.String(),
				Value:        "123.45",
				DisplayName:  strPtr("Balance"),
				Institution:  strPtr("Institution"),
				OfficialName: strPtr("Institution's Balance"),
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			assert.Equal(t, "EUR", res.Balance.Currency)
			assert.Equal(t, "CashChecking", res.Balance.Kind)
			assert.Equal(t, "Balance", *res.Balance.DisplayName)
			assert.Equal(t, "Institution", *res.Balance.Institution)
			assert.Equal(t, "Institution's Balance", *res.Balance.OfficialName)
		})

		t.Run("SuccessNoOptional", func(t *testing.T) {
			res, err := resolver.Mutation().CreateBalance(ctx, model.CreateBalanceInput{
				Currency: "RUB",
				Kind:     balance.CashPhysical.String(),
				Value:    "543.21",
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			assert.Equal(t, "RUB", res.Balance.Currency)
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
				case "EUR":
					assert.Equal(t, "CashChecking", b.Kind)
					assert.Equal(t, "Balance", *b.DisplayName)
					assert.Equal(t, "Institution", *b.Institution)
					assert.Equal(t, "Institution's Balance", *b.OfficialName)
					balanceID = b.ID
				case "RUB":
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
				BalanceID: balanceID,
				Value:     "234.56",
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
		})
	})

	t.Run("Balance_AllEntries", func(t *testing.T) {
		res, err := resolver.Balance().AllEntries(ctx, &model.Balance{ID: balanceID})
		require.NoError(t, err)
		require.Len(t, res, 2)

		assert.Equal(t, "234.56", res[0].Value)
		assert.Equal(t, "123.45", res[1].Value)
	})
}

func strPtr(s string) *string {
	return &s
}
