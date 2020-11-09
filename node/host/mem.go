package host

import (
	"fmt"
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

func getMemUsage() (*Usage, error) {
	mem, err := cmd.Exec("free", "--bytes")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(mem), "\n")
	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		values := utils.SplitSpaces(line)

		if values[0] == "Mem:" {
			total, err := strconv.ParseUint(values[1], 10, 64)
			if err != nil {
				return nil, err
			}
			used, err := strconv.ParseUint(values[2], 10, 64)
			if err != nil {
				return nil, err
			}
			usage := Usage{
				Total:      total,
				Used:       used,
				Percentage: utils.Round(float64(used) * 100 / float64(total)),
			}
			return &usage, nil
		}
	}

	return nil, fmt.Errorf("No memory found")
}
