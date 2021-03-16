// +build sqlite

package sqldb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteEngine struct{}

func (sqliteEngine) OpenDB(filePath string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filePath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func init() {
	registeredEngines["sqlite3"] = sqliteEngine{}
}
