package entry

import (
	"context"
	"time"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/shopspring/decimal"
)

// Repository describes the methods that must be implemented by a backend to access Entry entities.
type Repository interface {
	// List returns entries belonging to a specified balance.ID and matching the Filter.
	List(ctx context.Context, balanceID balance.ID, filter Filter) ([]*Entry, error)

	// CreateBatch inserts a batch of Entry entities from the map of currencies to values.
	CreateBatch(ctx context.Context, balanceID balance.ID, values map[currency.Currency]decimal.Decimal, validAt time.Time) error

	// DeleteAll removes all entries belonging to the specified balance.Balance.
	DeleteAll(ctx context.Context, balanceID balance.ID) error
}
