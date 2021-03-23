package prefs

//go:generate go run github.com/dmarkham/enumer -type=Key -output=key_gen.go -sql

type Key uint8

const (
	DefaultCurrency Key = iota
	SupportedCurrencies
)
