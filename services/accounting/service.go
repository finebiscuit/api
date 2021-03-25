package accounting

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
)

type Service interface {
	GetBalance(ctx context.Context, id balance.ID) (*balance.Balance, error)
	GetBalanceWithCurrentValue(ctx context.Context, id balance.ID) (*balance.WithCurrentValue, error)
	// GetHistoricalValuesPerMonth(ctx context.Context, balanceID balance.ID, cur currency.Currency, since, until time.Time) ([]*entry.Entry, error)
	ListBalancesWithCurrentValue(ctx context.Context) ([]*balance.WithCurrentValue, error)
	CreateBalance(ctx context.Context, b *balance.Balance, value decimal.Decimal) (balance.ValueMap, error)
	UpdateBalanceInfo(ctx context.Context, b *balance.Balance) error
	UpdateBalanceValue(ctx context.Context, balanceID balance.ID, value decimal.Decimal) (balance.ValueMap, error)
	DeleteBalance(ctx context.Context, balanceID balance.ID) error
}

type service struct {
	Tx TxFn
}

func NewService(tx TxFn) Service {
	return &service{Tx: tx}
}

func (s service) GetBalance(ctx context.Context, id balance.ID) (*balance.Balance, error) {
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

func (s service) GetBalanceWithCurrentValue(ctx context.Context, id balance.ID) (*balance.WithCurrentValue, error) {
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

func (s service) ListBalancesWithCurrentValue(ctx context.Context) ([]*balance.WithCurrentValue, error) {
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

func (s service) CreateBalance(ctx context.Context, b *balance.Balance, value decimal.Decimal) (balance.ValueMap, error) {
	if !b.Currency.IsACurrency() {
		return nil, fmt.Errorf("invalid or unsupported currency: %s", b.Currency)
	}

	var values balance.ValueMap
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		p, err := uow.Preferences().Get(ctx)
		if err != nil {
			return err
		}

		if !p.IsCurrencySupported(b.Currency) {
			return fmt.Errorf("currency not supported on this account (check preferences): %s", b.Currency)
		}

		values = map[currency.Currency]decimal.Decimal{
			b.Currency: value,
		}

		if b.Currency != p.DefaultCurrency {
			rate, err := uow.Forex().GetRate(ctx, b.Currency, p.DefaultCurrency)
			if err != nil {
				return err
			}
			values[p.DefaultCurrency] = values[b.Currency].Mul(rate)
		}

		if err := uow.Balances().Create(ctx, b); err != nil {
			return err
		}

		if err := uow.Entries().CreateBatch(ctx, b.ID, values, time.Now()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (s service) UpdateBalanceInfo(ctx context.Context, b *balance.Balance) error {
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

func (s service) UpdateBalanceValue(ctx context.Context, balanceID balance.ID, value decimal.Decimal) (balance.ValueMap, error) {
	var values balance.ValueMap
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		p, err := uow.Preferences().Get(ctx)
		if err != nil {
			return err
		}

		b, err := uow.Balances().GetWithCurrentValue(ctx, balanceID)
		if err != nil {
			return err
		}

		values = map[currency.Currency]decimal.Decimal{
			b.Currency: value,
		}

		if b.Currency != p.DefaultCurrency {
			rate, err := uow.Forex().GetRate(ctx, b.Currency, p.DefaultCurrency)
			if err != nil {
				return err
			}
			values[p.DefaultCurrency] = values[b.Currency].Mul(rate)
		}

		if err := uow.Entries().CreateBatch(ctx, balanceID, values, time.Now()); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (s service) DeleteBalance(ctx context.Context, balanceID balance.ID) error {
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
