package graph

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
	"github.com/finebiscuit/api/services/forex"
)

type AccountingService interface {
	ListBalances(ctx context.Context) ([]*balance.Balance, error)
	CreateBalance(ctx context.Context, b *balance.Balance, es []*entry.Entry) error
	UpdateBalance(ctx context.Context, b *balance.Balance) error
	AddEntries(ctx context.Context, es []*entry.Entry) error
	ListBalanceEntries(ctx context.Context, bID balance.ID) ([]*entry.Entry, error)
	FindBalance(ctx context.Context, id balance.ID, filter balance.Filter) (*balance.Balance, error)
	FindLatestEntryByCurrency(ctx context.Context, bID balance.ID, currency forex.Currency) (*entry.Entry, error)
}

type ForexService interface {
	GetRate(ctx context.Context, from, to forex.Currency) (decimal.Decimal, error)
}
