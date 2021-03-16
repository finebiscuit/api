package sqldb

import (
	"fmt"
	"strconv"
)

func parseID(s fmt.Stringer) uint {
	i, err := strconv.Atoi(s.String())
	if err != nil {
		return 0
	}
	return uint(i)
}

func uintToString(v uint) string {
	return strconv.Itoa(int(v))
}
