package accounting_test

import (
	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/services/accounting"
)

var _ graph.AccountingService = &accounting.Service{}
