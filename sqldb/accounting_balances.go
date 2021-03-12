package sqldb

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/finebiscuit/api/services/accounting/balance"
	"github.com/finebiscuit/api/services/forex"
	"github.com/finebiscuit/api/util"
)

type accountingBalancesRepository struct {
	db *gorm.DB
}

var _ balance.Repository = &accountingBalancesRepository{}

type accountingBalance struct {
	ID         balance.ID     `gorm:"size:20;unique;not null;check:length(id) >= 12"`
	Currency   forex.Currency `gorm:"size:8;not null;check:length(currency) >= 3"`
	Encryption string         `gorm:"not null;check:length(encryption) > 0"`
	Data       string         `gorm:"not null;check:length(encryption) > 0"`

	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func (accountingBalance) TableName() string { return "accounting_balances" }

func newAccountingBalance(b *balance.Balance) *accountingBalance {
	return &accountingBalance{
		ID:         b.ID,
		Currency:   b.Currency,
		Encryption: b.Encryption,
		Data:       b.Data,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  sql.NullTime{Time: b.UpdatedAt, Valid: !b.UpdatedAt.IsZero()},
	}
}

func (b *accountingBalance) toDomainEntity() *balance.Balance {
	return &balance.Balance{
		ID:       b.ID,
		Currency: b.Currency,
		EncryptedData: util.EncryptedData{
			Encryption: b.Encryption,
			Data:       b.Data,
		},
		HasTimestamps: util.HasTimestamps{
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt.Time,
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

func (r *accountingBalancesRepository) Get(ctx context.Context, id balance.ID, filter balance.Filter) (*balance.Balance, error) {
	var b accountingBalance
	if res := r.db.WithContext(ctx).First(&b, id); res.Error != nil {
		return nil, res.Error
	}
	return b.toDomainEntity(), nil
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

func (r *accountingBalancesRepository) Create(ctx context.Context, b *balance.Balance) error {
	ab := newAccountingBalance(b)
	if res := r.db.WithContext(ctx).Create(ab); res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *accountingBalancesRepository) Update(ctx context.Context, b *balance.Balance) error {
	panic("implement me")
}

func (r *accountingBalancesRepository) Delete(ctx context.Context, id balance.ID) error {
	panic("implement me")
}
