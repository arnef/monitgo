package utils

import "strings"

func SplitSpaces(val string) []string {
	var values []string
	for _, v := range strings.Split(val, " ") {
		trimed := strings.TrimSpace(v)
		if v != "" {
			values = append(values, trimed)
		}
	}
	return values
}
