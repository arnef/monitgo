package host

import (
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"

	"git.arnef.de/monitgo/log"
)

func getMemUsage() (map[string]Usage, error) {
	mem, err := cmd.Exec("free", "--bytes")
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	lines := strings.Split(string(mem), "\n")
	lines = lines[1 : len(lines)-1]
	usage := make(map[string]Usage)
	for _, line := range lines {
		values := utils.SplitSpaces(line)

		name := values[0][:len(values[0])-1]
		totalBytes, err := strconv.ParseUint(values[1], 10, 64)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		usedBytes, err := strconv.ParseUint(values[2], 10, 64)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		usage[name] = Usage{
			TotalBytes: totalBytes,
			UsedBytes:  usedBytes,
		}
	}

	return usage, nil
}
