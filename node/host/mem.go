package host

import (
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

type MemUsage struct {
	Name  string
	Total int
	Used  int
}

func getMemUsage() ([]MemUsage, error) {
	mem, err := cmd.Exec("free", "--mega")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(mem), "\n")
	lines = lines[1 : len(lines)-1]

	memory := make([]MemUsage, len(lines))

	for i, line := range lines {
		values := utils.SplitSpaces(line)

		total, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}
		used, err := strconv.Atoi(values[2])
		if err != nil {
			return nil, err
		}

		memory[i] = MemUsage{
			Name:  string(values[0][:len(values[0])-1]),
			Total: total,
			Used:  used,
		}
	}

	return memory, nil
}
