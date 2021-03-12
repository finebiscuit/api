package forex

type Currency string

func (c Currency) String() string {
	return string(c)
}
