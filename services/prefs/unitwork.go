package prefs

import (
	"context"
)

type TxFn func(ctx context.Context, fn func(ctx context.Context, uow UnitOfWork) error) error

type UnitOfWork interface {
	Preferences() Repository
}
