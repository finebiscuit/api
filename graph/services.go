package graph

import (
	"context"

	"github.com/finebiscuit/api/services/prefs"
	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
)

type AccountingService interface {
	GetBalance(ctx context.Context, id balance.ID) (*balance.Balance, error)
	GetBalanceWithCurrentValue(ctx context.Context, id balance.ID) (*balance.WithCurrentValue, error)
	ListBalancesWithCurrentValue(ctx context.Context) ([]*balance.WithCurrentValue, error)
	CreateBalance(ctx context.Context, b *balance.Balance, value decimal.Decimal) (balance.ValueMap, error)
	UpdateBalanceInfo(ctx context.Context, b *balance.Balance) error
	UpdateBalanceValue(ctx context.Context, balanceID balance.ID, value decimal.Decimal) (balance.ValueMap, error)
	DeleteBalance(ctx context.Context, balanceID balance.ID) error
}

type PreferencesService interface {
	GetPreferences(ctx context.Context) (*prefs.Preferences, error)
	SetPreferences(ctx context.Context, p *prefs.Preferences) error
}
