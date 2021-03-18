package entry

import (
	"context"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex"
)

// Repository describes the methods that must be implemented by a backend to access Entry entities.
type Repository interface {
	// Get returns a single Entry by ID.
	Get(ctx context.Context, balanceID balance.ID, entryID ID) (*Entry, error)

	// Find returns a single most recent Entry that matches the Filter.
	Find(ctx context.Context, balanceID balance.ID, filter Filter) (*Entry, error)

	// List returns entries belonging to a specified balance.ID and matching the Filter.
	List(ctx context.Context, balanceID balance.ID, filter Filter) ([]*Entry, error)

	// ListLatestPerBalanceAndCurrency returns all matching entries that belong to one of the specified balanceIDs.
	// The entries are returned in a double map, keyed by balance.ID on the first level and forex.Currency on the second.
	ListLatestPerBalanceAndCurrency(ctx context.Context, balanceIDs []balance.ID, filter Filter) (map[balance.ID]map[forex.Currency]*Entry, error)

	// Create creates a new Entry. The ID on the object is set to the last inserted ID.
	Create(ctx context.Context, e *Entry) error

	// Update updates an existing Entry, setting all fields.
	Update(ctx context.Context, e *Entry) error
}
