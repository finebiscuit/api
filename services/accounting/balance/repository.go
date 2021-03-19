package balance

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, id ID) (*Balance, error)
	GetWithCurrentValue(ctx context.Context, id ID) (*WithCurrentValue, error)
	List(ctx context.Context, filter Filter) ([]*Balance, error)
	ListWithCurrentValue(ctx context.Context, filter Filter) ([]*WithCurrentValue, error)
	Create(ctx context.Context, b *Balance) error
	Update(ctx context.Context, b *Balance) error
	Delete(ctx context.Context, id ID) error
}
