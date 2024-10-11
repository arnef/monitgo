package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func Free(val string) (uint64, uint64, error) {

	lines := strings.Split(string(val), "\n")

	var line string
	if len(lines) > 2 && len(lines[1]) > 0 {
		line = lines[1 : len(lines)-1][0]
	} else {
		return 0, 0, fmt.Errorf("invalid val: %s", val)
	}

	values := SplitSpaces(line)
	if len(values) < 3 {
		return 0, 0, fmt.Errorf("invalid val: %s", values)
	}

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
