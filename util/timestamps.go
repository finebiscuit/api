package util

import (
	"time"
)

type HasTimestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
