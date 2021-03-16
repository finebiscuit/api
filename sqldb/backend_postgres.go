// +build postgres

package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresEngine struct{}

func (postgresEngine) OpenDB(filePath string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(filePath))
}

func init() {
	registeredEngines["postgres"] = postgresEngine{}
}
