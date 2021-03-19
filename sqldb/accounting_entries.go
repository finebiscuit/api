package sqldb

import (
	"context"
	"time"

	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/accounting/entry"
	"github.com/finebiscuit/api/util"
)

type accountingEntriesRepository struct {
	db *gorm.DB
}

var _ entry.Repository = &accountingEntriesRepository{}

type accountingEntry struct {
	gorm.Model

	BalanceID uint
	Currency  currency.Currency `gorm:"size:8;not null;check:length(currency) >= 3"`
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
func (r *accountingEntriesRepository) List(ctx context.Context, balanceID balance.ID, filter entry.Filter) ([]*entry.Entry, error) {
	var aes []*accountingEntry

	q := filterEntry(r.db.WithContext(ctx), filter)
	if balanceID.Valid() {
		q = q.Where("balance_id = ?", balanceID)
	}
	if res := q.Order("valid_at desc").Find(&aes); res.Error != nil {
		return nil, res.Error
	}

	es := make([]*entry.Entry, 0, len(aes))
	for _, ae := range aes {
		es = append(es, ae.toDomainEntity())
	}
	return es, nil
}

func (r *accountingEntriesRepository) ListLatestPerBalanceAndCurrency(ctx context.Context, balanceIDs []balance.ID, filter entry.Filter) (map[balance.ID]map[currency.Currency]*entry.Entry, error) {
	result := make(map[balance.ID]map[currency.Currency]*entry.Entry)

	// TODO: improve, using multiple queries so very inefficient
	for _, bID := range balanceIDs {
		subQuery := filterEntry(r.db.Model(&accountingEntry{}), filter).Where("balance_id = ?", bID).Order("valid_at desc")
		q := r.db.WithContext(ctx).Table("(?) AS e", subQuery).Group("e.currency")

		var aes []*accountingEntry
		if res := q.Find(&aes); res.Error != nil {
			return nil, res.Error
		}

		result[bID] = make(map[currency.Currency]*entry.Entry)
		for _, ae := range aes {
			result[bID][ae.Currency] = ae.toDomainEntity()
		}
	}

	return result, nil
}

func (r *accountingEntriesRepository) CreateBatch(ctx context.Context, balanceID balance.ID, values map[currency.Currency]decimal.Decimal, validAt time.Time) error {
	for k, v := range values {
		e := &entry.Entry{
			BalanceID: balanceID,
			Currency:  k,
			Value:     v,
			ValidAt:   validAt,
		}
		ae := newAccountingEntry(e)
		if res := r.db.WithContext(ctx).Create(&ae); res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func (r *accountingEntriesRepository) DeleteAll(ctx context.Context, balanceID balance.ID) error {
	panic("implement me")
}

func filterEntry(q *gorm.DB, filter entry.Filter) *gorm.DB {
	params := make(map[string]interface{})
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
