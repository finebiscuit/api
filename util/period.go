package util

import (
	"fmt"
	"time"
)

type Period struct {
	Year  int
	Month time.Month
}

func NewPeriodFromTime(t time.Time) Period {
	return Period{Year: t.Year(), Month: t.Month()}
}

func (p Period) String() string {
	return fmt.Sprintf("%04d-%02d", p.Year, p.Month)
}

func (p Period) Time() time.Time {
	return time.Date(p.Year, p.Month, 1, 0, 0, 0, 0, time.UTC)
}

func (p Period) Range() (time.Time, time.Time) {
	since := p.Time()
	until := since.AddDate(0, 1, 0).Add(-1 * time.Millisecond)
	return since, until
}

func (p Period) Previous() Period {
	return NewPeriodFromTime(p.Time().AddDate(0, -1, 0))
}

func (p Period) Next() Period {
	return NewPeriodFromTime(p.Time().AddDate(0, 1, 0))
}
