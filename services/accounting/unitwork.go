package accounting

import (
	"context"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
)

type TxFn func(ctx context.Context, fn func(ctx context.Context, uow UnitOfWork) error) error

type UnitOfWork interface {
	Balances() balance.Repository
	Entries() entry.Repository
}
