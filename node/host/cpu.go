package host

import (
	"encoding/json"
	"fmt"
	"strconv"

	"git.arnef.de/monitgo/log"
	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/node/host/uptime"
	"git.arnef.de/monitgo/utils"
)

func getNormalizedLoad() (float64, error) {
	cpus, err := getCPUs()
	if err != nil || cpus == 0 {
		log.Debug(err, cpus)
		cpus = 4 // TODO get cpu count alternative
		// return nil, err
	}
	load, err := getLoad()
	if err != nil {
		log.Debug(err)
		return 0, err
	}
	return utils.Round(load[0] / float64(cpus)), nil
}

func getCPUs() (int, error) {
	lscpu, err := cmd.Exec("lscpu", "--json")
	if err != nil {
		log.Debug(err)
		return 0, err
	}
	var data map[string][]map[string]string
	err = json.Unmarshal(lscpu, &data)
	if err != nil {
		log.Debug(err)
		return 0, err
	}

	for _, line := range data["lscpu"] {
		if line["field"] == "CPU(s):" {
			val, err := strconv.Atoi(line["data"])
			if err != nil {
				log.Debug(err)
			}
			return val, err
		}
	}

	return 0, fmt.Errorf("could not get CPUs")
}

func getLoad() ([]float64, error) {
	uptimeOut, err := cmd.Exec("uptime")
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return uptime.ParseLoad(string(uptimeOut))
}
