package host

import (
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

type DiskUsage struct {
	Name  string
	Total int
	Used  int
}

func getDiskUsage() ([]DiskUsage, error) {
	df, err := cmd.Exec("df", "-BMB", "--output=source,size,used")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(df), "\n")
	lines = lines[1 : len(lines)-1]
	var usage []DiskUsage
	for _, line := range lines {
		values := utils.SplitSpaces(line)
		if values[0][0] == '/' {
			usage = append(usage, DiskUsage{
				Name:  values[0],
				Total: int(utils.MustParseMegabyte(values[1])),
				Used:  int(utils.MustParseMegabyte(values[2])),
			})
		}
	}

	return usage, nil
}
