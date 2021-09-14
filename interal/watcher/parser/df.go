package parser

import (
	"strconv"
	"strings"
)

func Df(val string) (uint64, uint64, error) {
	lines := strings.Split(val, "\n")
	lines = lines[1 : len(lines)-1]
	var total uint64 = 0
	var used uint64 = 0
	for _, line := range lines {
		values := SplitSpaces(line)
		if values[0][0] == '/' {
			totalBytes, err := strconv.ParseUint(values[1], 10, 64)
			if err != nil {
				return 0, 0, err
			}
			usedBytes, err := strconv.ParseUint(values[2], 10, 64)
			if err != nil {
				return 0, 0, err
			}
			total += totalBytes
			used += usedBytes

		}
	}
	return total, used, nil
}
