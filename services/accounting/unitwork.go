package accounting

import (
	"context"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
)

type UnitOfWorkStarter func(ctx context.Context) (UnitOfWork, Tx, error)

type Tx interface {
	Commit() error
	Rollback() error
}

type UnitOfWork interface {
	Balances() balance.Repository
	Entries() entry.Repository
}

func (s *Service) startUnitOfWork(ctx context.Context, fn func(uow UnitOfWork) error) error {
	uow, tx, err := s.StartUnitOfWork(ctx)
	if err != nil {
		return err
	}

	if err := fn(uow); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
