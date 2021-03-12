package model

import (
	"time"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
)

type Balance struct {
	ID         string `json:"id"`
	Currency   string `json:"currency"`
	Encryption string `json:"encryption"`
	Data       string `json:"data"`
}

func NewBalance(b *balance.Balance) *Balance {
	return &Balance{
		ID:         b.ID.String(),
		Currency:   b.Currency.String(),
		Encryption: b.Encryption,
		Data:       b.Data,
	}
}

type Entry struct {
	ID         string    `json:"id"`
	Currency   string    `json:"currency"`
	Encryption string    `json:"encryption"`
	Data       string    `json:"data"`
	ValidAt    time.Time `json:"validAt"`
}

func NewEntry(e *entry.Entry) *Entry {
	return &Entry{
		ID:         e.ID.String(),
		Currency:   e.Currency.String(),
		Encryption: e.Encryption,
		Data:       e.Data,
		ValidAt:    e.ValidAt,
	}
}
