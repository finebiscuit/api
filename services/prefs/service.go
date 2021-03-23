package prefs

import (
	"context"
	"fmt"
	"reflect"

	"github.com/finebiscuit/api/services/forex/currency"
)

type Service struct {
	Tx TxFn
}

func (s Service) GetPreferences(ctx context.Context) (*Preferences, error) {
	var p *Preferences

	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		var err error
		p, err = uow.Preferences().Get(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s Service) SetPreferences(ctx context.Context, p *Preferences) error {
	err := s.Tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		oldPrefs, err := uow.Preferences().Get(ctx)
		if err != nil {
			return err
		}

		changes, err := p.generateChanges(oldPrefs)
		if err != nil {
			return err
		}

		if err := uow.Preferences().Update(ctx, changes); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Preferences) generateChanges(before *Preferences) ([]Change, error) {
	changes := make([]Change, 0)

	if before.DefaultCurrency != p.DefaultCurrency {
		if !p.DefaultCurrency.IsACurrency() {
			return nil, fmt.Errorf("currency is invalid or unsupported: %s", p.DefaultCurrency.String())
		}
		changes = append(changes, Change{Key: DefaultCurrency, Value: p.DefaultCurrency})
	} else if !before.DefaultCurrency.IsACurrency() {
		return nil, fmt.Errorf("currency is invalid or unsupported: %s", before.DefaultCurrency.String())
	}

	beforeSuppMap := make(map[currency.Currency]struct{})
	for _, c := range before.SupportedCurrencies {
		beforeSuppMap[c] = struct{}{}
	}
	afterSuppMap := map[currency.Currency]struct{}{p.DefaultCurrency: {}}
	for _, c := range p.SupportedCurrencies {
		if !c.IsACurrency() {
			return nil, fmt.Errorf("currency is invalid or unsupported: %s", c.String())
		}
		afterSuppMap[c] = struct{}{}
	}
	if !reflect.DeepEqual(beforeSuppMap, afterSuppMap) {
		supp := make([]currency.Currency, 0, len(afterSuppMap))
		for c := range afterSuppMap {
			supp = append(supp, c)
		}
		p.SupportedCurrencies = supp
		changes = append(changes, Change{Key: SupportedCurrencies, Value: supp})
	}

	return changes, nil
}
