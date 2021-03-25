package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/services/accounting"
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/services/forex/dinero"
	"github.com/finebiscuit/api/services/prefs"
	"github.com/finebiscuit/api/services/projecting"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Accounting accounting.Service
	Forex      forex.Service
	Prefs      prefs.Service
	Projecting projecting.Service
}

type Backend interface {
	SupportedTypes() []string
	OpenAndPrepare(ctx context.Context, cfg *config.Config) error
	AccountingTxFn() accounting.TxFn
	PreferencesTxFn() prefs.TxFn
}

func NewResolver(cfg *config.Config, backends ...Backend) (*Resolver, error) {
	backend, ok := getBackend(backends, cfg.DBType)
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %q", cfg.DBType)
	}

	if err := backend.OpenAndPrepare(context.Background(), cfg); err != nil {
		return nil, err
	}

	prefsSvc := prefs.NewService(backend.PreferencesTxFn())
	accountingSvc := accounting.NewService(backend.AccountingTxFn())
	projectingSvc := projecting.NewService(accountingSvc)
	resolver := &Resolver{
		Forex:      forex.NewService(dinero.New("", 2*time.Hour)),
		Accounting: accountingSvc,
		Prefs:      prefsSvc,
		Projecting: projectingSvc,
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
