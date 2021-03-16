package sqldb

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/util"
)

type accountingEntriesRepository struct {
	db *gorm.DB
}

var _ entry.Repository = &accountingEntriesRepository{}

type accountingEntry struct {
	gorm.Model

	BalanceID uint
	Currency  forex.Currency `gorm:"size:8;not null;check:length(currency) >= 3"`
	Value     decimal.Decimal
	ValidAt   time.Time

	Balance accountingBalance
}

func (accountingEntry) TableName() string { return "accounting_entries" }

func newAccountingEntry(e *entry.Entry) *accountingEntry {
	return &accountingEntry{
		Model: gorm.Model{
			ID:        parseID(e.ID),
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		BalanceID: parseID(e.BalanceID),
		Currency:  e.Currency,
		Value:     e.Value,
		ValidAt:   e.ValidAt,
	}
}

func (e *accountingEntry) toDomainEntity() *entry.Entry {
	return &entry.Entry{
		ID:        entry.ParseID(uintToString(e.ID)),
		BalanceID: balance.ParseID(uintToString(e.BalanceID)),
		Currency:  e.Currency,
		Value:     e.Value,
		ValidAt:   e.ValidAt,
		HasTimestamps: util.HasTimestamps{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
	}
}

func (r *accountingEntriesRepository) Get(ctx context.Context, id entry.ID, filter entry.Filter) (*entry.Entry, error) {
	var e *accountingEntry

	q := filterEntry(r.db.WithContext(ctx), filter)
	if id.Valid() {
		q = q.Where("id = ?", id)
	}

	if res := q.Order("valid_at desc").First(&e); res.Error != nil {
		return nil, res.Error
	}

	return e.toDomainEntity(), nil
}

func (r *accountingEntriesRepository) List(ctx context.Context, filter entry.Filter) ([]*entry.Entry, error) {
	var aes []*accountingEntry

	q := filterEntry(r.db.WithContext(ctx), filter)
	if res := q.Order("valid_at desc").Find(&aes); res.Error != nil {
		return nil, res.Error
	}

	es := make([]*entry.Entry, 0, len(aes))
	for _, ae := range aes {
		es = append(es, ae.toDomainEntity())
	}
	return es, nil
}

func (r *accountingEntriesRepository) Create(ctx context.Context, e *entry.Entry) error {
	ae := newAccountingEntry(e)
	if res := r.db.WithContext(ctx).Create(&ae); res.Error != nil {
		return res.Error
	}
	e.ID = entry.ParseID(uintToString(ae.ID))
	return nil
}

func (r *accountingEntriesRepository) Update(ctx context.Context, e *entry.Entry) error {
	panic("implement me")
}

func filterEntry(q *gorm.DB, filter entry.Filter) *gorm.DB {
	params := make(map[string]interface{})
	if filter.BalanceIDs != nil {
		params["balance_id"] = filter.BalanceIDs
	}
	if filter.Currencies != nil {
		params["currency"] = filter.Currencies
	}

	q = q.Where(params)
	if filter.ValidAfter != nil {
		q = q.Where("valid_at > ?", *filter.ValidAfter)
	}
	if filter.ValidBefore != nil {
		q = q.Where("valid_at < ?", *filter.ValidBefore)
	}

	return q
}
