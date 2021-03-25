// +build postgres

package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/sqldb"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/finebiscuit/api/config"
)

func TestPostgres_Accounting(t *testing.T) {
	resolver := preparePostgres(t, "test_accounting")
	AccountingTests(t, context.Background(), resolver)
}

func TestPostgres_Preferences(t *testing.T) {
	resolver := preparePostgres(t, "test_preferences")
	PreferencesTests(t, context.Background(), resolver)
}

func TestPostgres_Projecting(t *testing.T) {
	resolver := preparePostgres(t, "test_projecting")
	ProjectingTests(t, context.Background(), resolver)
}

func postgresDSN(dbname string) string {
	return fmt.Sprintf("host=postgres port=5432 user=finebiscuit dbname=%s password=finebiscuit sslmode=disable", dbname)
}

func preparePostgres(t *testing.T, dbName string) *graph.Resolver {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cfg := &config.Config{
		DBType:   "postgres",
		DBSource: postgresDSN(dbName),
	}

	db, err := gorm.Open(postgres.Open(postgresDSN("postgres")))
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
	require.NoError(t, db.Error)

	db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	require.NoError(t, db.Error)

	resolver, err := graph.NewResolver(cfg, sqldb.NewBackend(forexMock{}))
	require.NoError(t, err)

	return resolver
}
