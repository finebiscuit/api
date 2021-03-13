// +build sqlite

package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/finebiscuit/api/config"
)

func TestSqlite_Accounting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	ctx := context.Background()
	cfg := &config.Config{
		DBType:   "sqlite3",
		DBSource: "accounting.test.db",
	}

	t.Cleanup(func() {
		os.Remove(cfg.DBSource)
	})

	accountingTests(t, ctx, cfg)
}
