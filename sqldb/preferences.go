package sqldb

import (
	"context"
	"fmt"
	"strings"

	"github.com/finebiscuit/api/services/forex/currency"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/finebiscuit/api/services/prefs"
)

type preferencesRepository struct {
	db *gorm.DB
}

var _ prefs.Repository = &preferencesRepository{}

type preference struct {
	Key   prefs.Key `gorm:"primarykey"`
	Value string
}

func (preference) TableName() string { return "preferences" }

func (r preferencesRepository) Get(ctx context.Context) (*prefs.Preferences, error) {
	var ps []preference
	if res := r.db.WithContext(ctx).Find(&ps); res.Error != nil {
		return nil, res.Error
	}

	var result prefs.Preferences
	for _, p := range ps {
		switch p.Key {
		case prefs.DefaultCurrency:
			cur, _ := currency.CurrencyString(p.Value)
			result.DefaultCurrency = cur
		case prefs.SupportedCurrencies:
			ss := strings.Split(p.Value, ",")
			result.SupportedCurrencies = make([]currency.Currency, 0, len(ss))
			for _, s := range ss {
				cur, _ := currency.CurrencyString(s)
				result.SupportedCurrencies = append(result.SupportedCurrencies, cur)
			}
		}
	}
	return &result, nil
}

func (r preferencesRepository) Update(ctx context.Context, changes []prefs.Change) error {
	for _, ch := range changes {
		var value string
		switch v := ch.Value.(type) {
		case string:
			value = v
		case currency.Currency:
			value = v.String()
		case []currency.Currency:
			ss := make([]string, 0, len(v))
			for _, x := range v {
				ss = append(ss, x.String())
			}
			value = strings.Join(ss, ",")
		default:
			return fmt.Errorf("invalid preference value type: %T", ch.Value)
		}
		p := preference{Key: ch.Key, Value: value}
		if res := r.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&p); res.Error != nil {
			return res.Error
		}
	}
	return nil
}
