package utils

import (
	"strconv"
	"strings"
)

func AddCanonicalization(a float64, b float64) string {
	return strings.Join([]string{strconv.FormatFloat(a, 'f', -1, 64), strconv.FormatFloat(b, 'f', -1, 64)}, ",")
}
