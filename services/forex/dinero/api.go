package dinero

import (
	"context"
	"time"

	"github.com/mattevans/dinero"
	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/forex/currency"
)

type API struct {
	client *dinero.Client
}

func New(appID string, cacheExpiry time.Duration) *API {
	client := dinero.NewClient(appID, "USD", cacheExpiry)
	return &API{client: client}
}

func (api *API) GetRate(ctx context.Context, from, to currency.Currency) (decimal.Decimal, error) {
	api.client.Rates.SetBaseCurrency(string(from))
	rate, err := api.client.Rates.Get(string(to))
	if err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromFloat(*rate), nil
}
