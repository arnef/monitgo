package uptime

import (
	"fmt"
	"strconv"
	"strings"

	"git.arnef.de/monitgo/log"
)

func ParseLoad(val string) ([]float64, error) {
	load := strings.Split(val, "load average: ")
	loads := strings.Split(load[len(load)-1], ", ")
	if len(loads) < 3 {
		return nil, fmt.Errorf("could not parse \"%s\"", val)
	}
	vals := make([]float64, 3)
	for i := 0; i < 3; i++ {
		cleanLoad := strings.TrimSpace(strings.Split(loads[i], " ")[0])
		cleanLoad = strings.Replace(cleanLoad, ",", ".", 1)
		parsed, err := strconv.ParseFloat(cleanLoad, 64)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		vals[i] = parsed * 100
	}
	return vals, nil
}
