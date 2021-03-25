// +build sqlite

package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/sqldb"
	"github.com/stretchr/testify/require"
)

func TestSqlite_Accounting(t *testing.T) {
	resolver := prepareSqlite(t, "accounting.test.db")
	AccountingTests(t, context.Background(), resolver)
}

func TestSqlite_Preferences(t *testing.T) {
	resolver := prepareSqlite(t, "preferences.test.db")
	PreferencesTests(t, context.Background(), resolver)
}

func TestSqlite_Projecting(t *testing.T) {
	resolver := prepareSqlite(t, "projecting.test.db")
	ProjectingTests(t, context.Background(), resolver)
}

func prepareSqlite(t *testing.T, dbFile string) *graph.Resolver {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cfg := &config.Config{
		DBType:   "sqlite3",
		DBSource: dbFile,
	}

	t.Cleanup(func() {
		os.Remove(cfg.DBSource)
	})

	resolver, err := graph.NewResolver(cfg, sqldb.NewBackend(forexMock{}))
	require.NoError(t, err)

	return resolver
}
