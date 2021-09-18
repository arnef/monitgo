package parser

import (
	"strconv"
	"strings"
)

func Free(val string) (uint64, uint64, error) {

	lines := strings.Split(string(val), "\n")
	line := lines[1 : len(lines)-1][0]

	values := SplitSpaces(line)

	totalBytes, err := strconv.ParseUint(values[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	usedBytes, err := strconv.ParseUint(values[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return totalBytes, usedBytes, nil
}
