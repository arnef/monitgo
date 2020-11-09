package host

import (
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

func getDiskUsage() (*Usage, error) {
	df, err := cmd.Exec("df", "--output=source,size,used")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(df), "\n")
	lines = lines[1 : len(lines)-1]
	usage := Usage{}
	for _, line := range lines {
		values := utils.SplitSpaces(line)
		if values[0][0] == '/' {
			total, err := strconv.ParseUint(values[1], 10, 64)
			if err != nil {
				return nil, err
			}
			used, err := strconv.ParseUint(values[2], 10, 64)
			if err != nil {
				return nil, err
			}
			usage.Used += used
			usage.Total += total
		}
	}

	usage.Percentage = utils.Round(float64(usage.Used) * 100 / float64(usage.Total))

	return &usage, nil
}
