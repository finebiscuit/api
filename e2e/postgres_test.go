// +build postgres

package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/config"
)

func TestPostgres_Accounting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx := context.Background()
	dbName := "test_accounting"
	cfg := &config.Config{
		DBType:   "postgres",
		DBSource: postgresDSN(dbName),
		JWTKey:   "keyboard-cat",
	}

	preparePostgres(t, dbName)
	accountingTests(t, ctx, cfg)
}

func postgresDSN(dbname string) string {
	return fmt.Sprintf("host=postgres port=5432 user=finebiscuit dbname=%s password=finebiscuit sslmode=disable", dbname)
}

func preparePostgres(t *testing.T, dbname string) {
	db, err := gorm.Open(postgres.Open(postgresDSN("postgres")))
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbname))
	require.NoError(t, db.Error)

	db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
	require.NoError(t, db.Error)
}
