package sqldb

import (
	"context"

	"gorm.io/gorm"

	"github.com/finebiscuit/api/services/accounting"
	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
)

type UnitOfWork struct {
	tx *gorm.DB
}

func newUnitOfWork(ctx context.Context, db *gorm.DB) (*UnitOfWork, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	uow := &UnitOfWork{tx: tx}
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

func (b *Backend) AccountingTxFn() accounting.TxFn {
	return func(ctx context.Context, fn func(ctx context.Context, uow accounting.UnitOfWork) error) error {
		uow, err := newUnitOfWork(ctx, b.db)
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
