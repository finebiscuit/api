package model

import (
	"time"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
)

type Balance struct {
	ID           string  `json:"id"`
	Currency     string  `json:"currency"`
	Kind         string  `json:"kind"`
	DisplayName  *string `json:"displayName"`
	OfficialName *string `json:"officialName"`
	Institution  *string `json:"institution"`
}

func NewBalance(b *balance.Balance) *Balance {
	return &Balance{
		ID:           string(b.ID),
		Currency:     b.Currency.String(),
		Kind:         b.Type.String(),
		DisplayName:  &b.DisplayName,
		OfficialName: &b.OfficialName,
		Institution:  &b.Institution,
	}
}

type Entry struct {
	ID       string    `json:"id"`
	Currency string    `json:"currency"`
	Value    string    `json:"value"`
	ValidAt  time.Time `json:"validAt"`
}

func NewEntry(e *entry.Entry) *Entry {
	return &Entry{
		ID:       string(e.ID),
		Currency: e.Currency.String(),
		Value:    e.Value.String(),
		ValidAt:  e.ValidAt,
	}
}
