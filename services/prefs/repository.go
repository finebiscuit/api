package prefs

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context) (*Preferences, error)
	Update(ctx context.Context, changes []Change) error
}
