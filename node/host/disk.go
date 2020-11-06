package host

import (
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

func getDiskUsage() (*Usage, error) {
	df, err := cmd.Exec("df", "-BMB", "--output=source,size,used")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(df), "\n")
	lines = lines[1 : len(lines)-1]
	usage := Usage{}
	for _, line := range lines {
		values := utils.SplitSpaces(line)
		if values[0][0] == '/' {
			usage.Used += utils.MustParseMegabyte(values[2])
			usage.Total += utils.MustParseMegabyte(values[1])
		}
	}

	usage.Percentage = utils.Round(usage.Used * 100 / usage.Total)

	return &usage, nil
}
