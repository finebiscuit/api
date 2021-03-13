package accounting

import (
	"context"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
	"github.com/finebiscuit/api/services/forex"
)

type Service struct {
	Tx TxFn
}

func (s Service) ListBalances(ctx context.Context) ([]*balance.Balance, error) {
	var bals []*balance.Balance

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		bals, err = uow.Balances().List(ctx, balance.Filter{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bals, nil
}

func (s Service) CreateBalance(ctx context.Context, b *balance.Balance, es []*entry.Entry) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Balances().Create(ctx, b); err != nil {
			return err
		}

		for _, e := range es {
			if err := uow.Entries().Create(ctx, e); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) UpdateBalance(ctx context.Context, b *balance.Balance) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Balances().Update(ctx, b); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) AddEntries(ctx context.Context, es []*entry.Entry) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		for _, e := range es {
			if err := uow.Entries().Create(ctx, e); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) ListBalanceEntries(ctx context.Context, bID balance.ID) ([]*entry.Entry, error) {
	var entries []*entry.Entry

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		filter := entry.Filter{BalanceIDs: []balance.ID{bID}}
		entries, err = uow.Entries().List(ctx, filter)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (s Service) FindBalance(ctx context.Context, id balance.ID, filter balance.Filter) (*balance.Balance, error) {
	var b *balance.Balance

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		b, err = uow.Balances().Get(ctx, id, filter)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s Service) FindLatestEntryByCurrency(ctx context.Context, bID balance.ID, currency forex.Currency) (*entry.Entry, error) {
	var e *entry.Entry

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		filter := entry.Filter{
			BalanceIDs: []balance.ID{bID},
			Currencies: []forex.Currency{currency},
		}
		e, err = uow.Entries().Get(ctx, entry.NoID, filter)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return e, nil
}
