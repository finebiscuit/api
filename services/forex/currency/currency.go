package currency

import (
	"strings"
)

//go:generate go run github.com/dmarkham/enumer -type=Currency -trimprefix Currency -output=currency_gen.go -gqlgen -sql -json -text

type Currency uint8

const (
	XXX Currency = iota + 1
	XTS

	// G10 currencies https://en.wikipedia.org/wiki/G10_currencies.
	AUD
	CAD
	EUR
	JPY
	NZD
	NOK
	GBP
	SEK
	CHF
	USD

	// Other supported currencies
	RUB
	TRY
	DKK
	PLN
	HUF
	CZK
	ILS
	AED
	RON
	BGN
	RSD
	UAH

	// Precious metals
	XAG
	XAU
	XPT
	XPD

	// Cryptocurrency
	BCH
	BTC
	ETH
	LTC
	XRP
)

func New(s string) Currency {
	c, err := CurrencyString(strings.ToUpper(s))
	if err != nil {
		return XXX
	}
	return c
}
