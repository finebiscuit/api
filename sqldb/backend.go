package sqldb

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/finebiscuit/api/config"
)

type Backend struct {
	db *gorm.DB
}

func NewBackend() *Backend {
	return &Backend{}
}

func (b *Backend) SupportedTypes() []string {
	ts := make([]string, 0, len(registeredEngines))
	for k := range registeredEngines {
		ts = append(ts, k)
	}
	return ts
}

func (b *Backend) OpenAndPrepare(ctx context.Context, cfg *config.Config) error {
	eng, ok := registeredEngines[cfg.DBType]
	if !ok {
		return fmt.Errorf("unsupported database type: %q", cfg.DBType)
	}

	db, err := eng.OpenDB(cfg.DBSource)
	if err != nil {
		return err
	}

	if err := db.WithContext(ctx).AutoMigrate(
		&accountingBalance{},
		&accountingEntry{},
	); err != nil {
		return err
	}

	b.db = db
	return nil
}

type engine interface {
	OpenDB(dbSource string) (*gorm.DB, error)
}

var registeredEngines = map[string]engine{}
