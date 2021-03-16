package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/services/accounting"
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/services/forex/dinero"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Accounting AccountingService
	Forex      ForexService
}

type Backend interface {
	SupportedTypes() []string
	OpenAndPrepare(ctx context.Context, cfg *config.Config) error
	AccountingTxFn() accounting.TxFn
}

func NewResolver(cfg *config.Config, backends ...Backend) (*Resolver, error) {
	backend, ok := getBackend(backends, cfg.DBType)
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %q", cfg.DBType)
	}

	if err := backend.OpenAndPrepare(context.Background(), cfg); err != nil {
		return nil, err
	}

	accountingService := &accounting.Service{Tx: backend.AccountingTxFn()}
	resolver := &Resolver{
		Forex:      forex.NewService(dinero.New("", 2*time.Hour)),
		Accounting: accountingService,
	}
	return resolver, nil
}

func getBackend(backends []Backend, dbType string) (Backend, bool) {
	for _, b := range backends {
		ts := b.SupportedTypes()
		for _, t := range ts {
			if dbType == t {
				return b, true
			}
		}
	}
	return nil, false
}
