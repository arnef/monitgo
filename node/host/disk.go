package host

import (
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

func getDiskUsage() (map[string]Usage, error) {
	df, err := cmd.Exec("df", "--output=source,size,used")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(df), "\n")
	lines = lines[1 : len(lines)-1]
	usage := make(map[string]Usage)
	for _, line := range lines {
		values := utils.SplitSpaces(line)
		if values[0][0] == '/' {
			totalBytes, err := strconv.ParseUint(values[1], 10, 64)
			if err != nil {
				return nil, err
			}
			usedBytes, err := strconv.ParseUint(values[2], 10, 64)
			if err != nil {
				return nil, err
			}
			usage[values[0]] = Usage{
				TotalBytes: totalBytes,
				UsedBytes:  usedBytes,
			}
		}
	}

	return usage, nil
}
