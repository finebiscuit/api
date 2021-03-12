package graph

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/services/accounting"
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/services/forex/dinero"
	"github.com/finebiscuit/api/sqldb"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Accounting AccountingService
	Forex      ForexService
}

func NewResolver(cfg *config.Config) (*Resolver, error) {
	var (
		accountingSUOW accounting.UnitOfWorkStarter
	)

	switch cfg.DBType {
	case "sqlite3", "postgres":
		db, err := openSQLDatabase(cfg.DBType, cfg.DBSource)
		if err != nil {
			return nil, err
		}
		accountingSUOW = sqldb.Accounting(db)

		if err := sqldb.AutoMigrate(context.Background(), db); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %q", cfg.DBType)
	}

	accountingService := &accounting.Service{
		StartUnitOfWork: accountingSUOW,
	}

	resolver := &Resolver{
		Forex:      forex.NewService(dinero.New("", 2*time.Hour)),
		Accounting: accountingService,
	}
	return resolver, nil
}

func openSQLDatabase(dbType, dbSource string) (*gorm.DB, error) {
	switch dbType {
	case "sqlite3":
		return gorm.Open(sqlite.Open(dbSource), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	case "postgres":
		return gorm.Open(postgres.Open(dbSource))
	default:
		return nil, fmt.Errorf("unsupported database type: %q", dbType)
	}
}
