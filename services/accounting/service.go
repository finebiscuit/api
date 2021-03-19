package accounting

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
)

type Service struct {
	Tx TxFn
}

func (s Service) GetBalance(ctx context.Context, id balance.ID) (*balance.Balance, error) {
	var b *balance.Balance

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		b, err = uow.Balances().Get(ctx, id)
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

func (s Service) GetBalanceWithCurrentValue(ctx context.Context, id balance.ID) (*balance.WithCurrentValue, error) {
	var b *balance.WithCurrentValue

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		b, err = uow.Balances().GetWithCurrentValue(ctx, id)
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

func (s Service) ListBalancesWithCurrentValue(ctx context.Context) ([]*balance.WithCurrentValue, error) {
	var bals []*balance.WithCurrentValue

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		bals, err = uow.Balances().ListWithCurrentValue(ctx, balance.Filter{})
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

func (s Service) CreateBalance(ctx context.Context, b *balance.Balance, values map[currency.Currency]decimal.Decimal) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Balances().Create(ctx, b); err != nil {
			return err
		}

		if err := uow.Entries().CreateBatch(ctx, b.ID, values, time.Now()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) UpdateBalanceInfo(ctx context.Context, b *balance.Balance) error {
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

func (s Service) UpdateBalanceValue(ctx context.Context, balanceID balance.ID, values map[currency.Currency]decimal.Decimal) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Entries().CreateBatch(ctx, balanceID, values, time.Now()); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) DeleteBalance(ctx context.Context, balanceID balance.ID) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Entries().DeleteAll(ctx, balanceID); err != nil {
			return err
		}
		if err := uow.Balances().Delete(ctx, balanceID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
