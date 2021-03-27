package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/graph/model"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/finebiscuit/api/services/prefs"
)

func ProjectingTests(t *testing.T, ctx context.Context, resolver *graph.Resolver) {
	err := resolver.Prefs.SetPreferences(ctx, &prefs.Preferences{
		DefaultCurrency:     currency.EUR,
		SupportedCurrencies: []currency.Currency{currency.EUR},
	})
	require.NoError(t, err)

	b := balance.New(currency.EUR, balance.CashChecking, balance.Optional{
		EstimatedMonthlyValueChange: decimal.NewFromInt(100),
		EstimatedMonthlyGrowthRate:  decimal.NewFromFloat(1.02),
	})
	valMap, err := resolver.Accounting.CreateBalance(ctx, b, decimal.NewFromInt(100))
	require.NoError(t, err)
	bwcv := &balance.WithCurrentValue{Balance: *b, CurrentValue: valMap, ValidAt: time.Now()}

	t.Run("Balance_ProjectedValues", func(t *testing.T) {
		c := currency.EUR
		year, month, _ := time.Now().AddDate(0, 1, 0).Date()
		values, err := resolver.Balance().ProjectedValues(ctx, model.NewBalance(bwcv), 6, &c)
		require.NoError(t, err)

		since := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		expected := []*model.BalanceValue{
			{
				Currency: currency.EUR,
				Value:    decimal.NewFromFloat(202), // 100 * 1.02 + 100
				ValidAt:  since,
			},
			{
				Currency: currency.EUR,
				Value:    decimal.NewFromFloat(306.04), // 202 * 1.02 + 100
				ValidAt:  since.AddDate(0, 1, 0),
			},
			{
				Currency: currency.EUR,
				Value:    decimal.NewFromFloat(412.16), // 306.04 * 1.02 + 100
				ValidAt:  since.AddDate(0, 2, 0),
			},
			{
				Currency: currency.EUR,
				Value:    decimal.NewFromFloat(520.4), // 412.16 * 1.02 + 100
				ValidAt:  since.AddDate(0, 3, 0),
			},
			{
				Currency: currency.EUR,
				Value:    decimal.NewFromFloat(630.81), // 520.4 * 1.02 + 100
				ValidAt:  since.AddDate(0, 4, 0),
			},
			{
				Currency: currency.EUR,
				Value:    decimal.NewFromFloat(743.43), // 630.81 * 1.02 + 100
				ValidAt:  since.AddDate(0, 5, 0),
			},
		}

		if !cmp.Equal(expected, values) {
			t.Errorf("Not equal:\n%s", cmp.Diff(expected, values))
		}
	})
}
