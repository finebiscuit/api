package sqldb

import (
	"context"

	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/services/prefs"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/services/accounting"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
)

type UnitOfWork struct {
	tx    *gorm.DB
	forex forex.Service
}

func newUnitOfWork(ctx context.Context, b *Backend) (*UnitOfWork, error) {
	tx := b.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	uow := &UnitOfWork{tx: tx, forex: b.Forex}
	return uow, nil
}

func (uow *UnitOfWork) Commit() error {
	return uow.tx.Commit().Error
}

func (uow *UnitOfWork) Rollback() error {
	return uow.tx.Rollback().Error
}

func (uow *UnitOfWork) Balances() balance.Repository {
	return &accountingBalancesRepository{db: uow.tx}
}

func (uow *UnitOfWork) Entries() entry.Repository {
	return &accountingEntriesRepository{db: uow.tx}
}

func (uow *UnitOfWork) Preferences() prefs.Repository {
	return &preferencesRepository{db: uow.tx}
}

func (uow *UnitOfWork) Forex() forex.Repository {
	return uow.forex
}

func (b *Backend) AccountingTxFn() accounting.TxFn {
	return func(ctx context.Context, fn func(ctx context.Context, uow accounting.UnitOfWork) error) error {
		uow, err := newUnitOfWork(ctx, b)
		if err != nil {
			return err
		}
		if err := fn(ctx, uow); err != nil {
			uow.Rollback()
			return err
		} else {
			if err := uow.Commit(); err != nil {
				uow.Rollback()
				return err
			}
			return nil
		}
	}
}

func (b *Backend) PreferencesTxFn() prefs.TxFn {
	return func(ctx context.Context, fn func(ctx context.Context, uow prefs.UnitOfWork) error) error {
		uow, err := newUnitOfWork(ctx, b)
		if err != nil {
			return err
		}
		if err := fn(ctx, uow); err != nil {
			uow.Rollback()
			return err
		} else {
			if err := uow.Commit(); err != nil {
				uow.Rollback()
				return err
			}
			return nil
		}
	}
}
