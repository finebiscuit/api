package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/graph/model"
)

func accountingTests(t *testing.T, ctx context.Context, cfg *config.Config) {
	resolver, err := graph.NewResolver(cfg)
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
				Currency:   "EUR",
				Encryption: "none",
				Data:       "type=Cash",
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			assert.Equal(t, "EUR", res.Balance.Currency)
			assert.Equal(t, "none", res.Balance.Encryption)
			assert.Equal(t, "type=Cash", res.Balance.Data)
		})
	})

	t.Run("Query_Balances", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			res, err := resolver.Query().Balances(ctx)
			require.NoError(t, err)
			require.Len(t, res, 1)

			b := res[0]
			assert.Equal(t, "EUR", b.Currency)
			assert.Equal(t, "none", b.Encryption)
			assert.Equal(t, "type=Cash", b.Data)

			balanceID = b.ID
		})
	})

	t.Run("Mutation_CreateEntries", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			res, err := resolver.Mutation().AddEntries(ctx, model.AddEntriesInput{
				BalanceID: balanceID,
				Entries: []*model.EntryInput{
					{
						Currency:   "EUR",
						Encryption: "none",
						Data:       "value=123.45",
					},
					{
						Currency:   "USD",
						Encryption: "none",
						Data:       "value=234.56",
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, res.Balance)
			assert.Equal(t, balanceID, res.Balance.ID)
		})
	})

	t.Run("Balance_AllEntries", func(t *testing.T) {
		res, err := resolver.Balance().AllEntries(ctx, &model.Balance{ID: balanceID})
		require.NoError(t, err)
		require.Len(t, res, 2)

		for _, e := range res {
			switch e.Currency {
			case "EUR":
				assert.Equal(t, "none", e.Encryption)
				assert.Equal(t, "value=123.45", e.Data)
			case "USD":
				assert.Equal(t, "none", e.Encryption)
				assert.Equal(t, "value=234.56", e.Data)
			}
		}
	})
}
