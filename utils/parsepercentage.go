package utils

import (
	"strconv"
	"strings"
)

// ParsePercentage removes % from string and returns the value as float
func ParsePercentage(val string) (float64, error) {

	return strconv.ParseFloat(strings.Replace(val, "%", "", 1), 64)
}

// MustParsePercentage parsed and panic on failure
func MustParsePercentage(val string) float64 {
	res, err := ParsePercentage(val)
	if err != nil {
		panic(err)
	}
	return res
}
