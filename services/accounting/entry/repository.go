package entry

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, id ID, filter Filter) (*Entry, error)
	List(ctx context.Context, filter Filter) ([]*Entry, error)
	Create(ctx context.Context, e *Entry) error
	Update(ctx context.Context, e *Entry) error
}
