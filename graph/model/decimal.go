package model

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shopspring/decimal"
)

func MarshalDecimal(d decimal.Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(d.String()))
	})
}

func UnmarshalDecimal(v interface{}) (decimal.Decimal, error) {
	s, ok := v.(string)
	if !ok {
		return decimal.Zero, fmt.Errorf("decimal must be a string")
	}
	return decimal.NewFromString(s)
}
