package sqldb

import (
	"context"
	"database/sql"

	"github.com/finebiscuit/api/services/accounting/entry"
	"github.com/finebiscuit/api/services/forex/currency"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/util"
)

type accountingBalancesRepository struct {
	db *gorm.DB
}

var _ balance.Repository = &accountingBalancesRepository{}

type accountingBalance struct {
	gorm.Model

	Currency currency.Currency `gorm:"size:8;not null;check:length(currency) >= 3"`
	Type     balance.Type

	DisplayName  sql.NullString
	OfficialName sql.NullString
	Institution  sql.NullString

	EstimatedMonthlyGrowthRate  decimal.NullDecimal
	EstimatedMonthlyValueChange decimal.NullDecimal
}

func (accountingBalance) TableName() string { return "accounting_balances" }

func newAccountingBalance(b *balance.Balance) *accountingBalance {
	return &accountingBalance{
		Model: gorm.Model{
			ID:        parseID(b.ID),
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},

		Currency: b.Currency,
		Type:     b.Type,

		DisplayName:  sql.NullString{String: b.DisplayName, Valid: b.DisplayName != ""},
		OfficialName: sql.NullString{String: b.OfficialName, Valid: b.OfficialName != ""},
		Institution:  sql.NullString{String: b.Institution, Valid: b.Institution != ""},

		EstimatedMonthlyValueChange: decimal.NullDecimal{Decimal: b.EstimatedMonthlyValueChange, Valid: !b.EstimatedMonthlyValueChange.IsZero()},
		EstimatedMonthlyGrowthRate:  decimal.NullDecimal{Decimal: b.EstimatedMonthlyGrowthRate, Valid: !b.EstimatedMonthlyGrowthRate.IsZero()},
	}
}

func (b *accountingBalance) toDomainEntity() *balance.Balance {
	return &balance.Balance{
		ID:       balance.ParseID(uintToString(b.ID)),
		Currency: b.Currency,
		Type:     b.Type,
		Optional: balance.Optional{
			DisplayName:                 b.DisplayName.String,
			OfficialName:                b.OfficialName.String,
			Institution:                 b.Institution.String,
			EstimatedMonthlyGrowthRate:  b.EstimatedMonthlyGrowthRate.Decimal,
			EstimatedMonthlyValueChange: b.EstimatedMonthlyValueChange.Decimal,
		},
		HasTimestamps: util.HasTimestamps{
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		},
	}
}

type accountingBalanceList []*accountingBalance

func (l accountingBalanceList) toDomainEntity() []*balance.Balance {
	bs := make([]*balance.Balance, 0, len(l))
	for _, b := range l {
		bs = append(bs, b.toDomainEntity())
	}
	return bs
}

func (r *accountingBalancesRepository) Get(ctx context.Context, id balance.ID) (*balance.Balance, error) {
	var b accountingBalance
	if res := r.db.WithContext(ctx).First(&b, id); res.Error != nil {
		return nil, res.Error
	}
	return b.toDomainEntity(), nil
}

func (r *accountingBalancesRepository) GetWithCurrentValue(ctx context.Context, id balance.ID) (*balance.WithCurrentValue, error) {
	b, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	entryRepo := &accountingEntriesRepository{db: r.db}
	m, err := entryRepo.ListLatestPerBalanceAndCurrency(ctx, []balance.ID{id}, entry.Filter{})
	if err != nil {
		return nil, err
	}

	valMap := make(map[currency.Currency]decimal.Decimal)
	for cur, e := range m[id] {
		valMap[cur] = e.Value
	}

	bwcv := &balance.WithCurrentValue{
		Balance:      *b,
		CurrentValue: valMap,
	}
	return bwcv, nil
}

func (r *accountingBalancesRepository) List(ctx context.Context, filter balance.Filter) ([]*balance.Balance, error) {
	var bs accountingBalanceList

	params := make(map[string]interface{})
	if filter.Currencies != nil {
		params["currency"] = filter.Currencies
	}

	if res := r.db.WithContext(ctx).Find(&bs); res.Error != nil {
		return nil, res.Error
	}
	return bs.toDomainEntity(), nil
}

func (r *accountingBalancesRepository) ListWithCurrentValue(ctx context.Context, filter balance.Filter) ([]*balance.WithCurrentValue, error) {
	bs, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	ids := make([]balance.ID, 0, len(bs))
	for _, b := range bs {
		ids = append(ids, b.ID)
	}

	entryRepo := &accountingEntriesRepository{db: r.db}
	m, err := entryRepo.ListLatestPerBalanceAndCurrency(ctx, ids, entry.Filter{})
	if err != nil {
		return nil, err
	}

	bwcvs := make([]*balance.WithCurrentValue, 0, len(bs))
	for _, b := range bs {
		valMap := make(map[currency.Currency]decimal.Decimal)
		for cur, e := range m[b.ID] {
			valMap[cur] = e.Value
		}

		bwcvs = append(bwcvs, &balance.WithCurrentValue{
			Balance:      *b,
			CurrentValue: valMap,
		})
	}
	return bwcvs, nil
}

func (r *accountingBalancesRepository) Create(ctx context.Context, b *balance.Balance) error {
	ab := newAccountingBalance(b)
	if res := r.db.WithContext(ctx).Create(&ab); res.Error != nil {
		return res.Error
	}
	b.ID = balance.ParseID(uintToString(ab.ID))
	return nil
}

func (r *accountingBalancesRepository) Update(ctx context.Context, b *balance.Balance) error {
	panic("implement me")
}

func (r *accountingBalancesRepository) Delete(ctx context.Context, id balance.ID) error {
	panic("implement me")
}
