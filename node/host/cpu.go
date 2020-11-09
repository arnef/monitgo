package host

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
	"git.arnef.de/monitgo/utils"
)

func getNormalizedLoad() ([]float64, error) {
	cpus, err := getCPUs()
	if err != nil {
		return nil, err
	}
	load, err := getLoad()
	if err != nil {
		return nil, err
	}
	for i, l := range load {
		load[i] = utils.Round(l / float64(cpus))
		// fmt.Println(l / cpus)
	}
	return load, nil
}

func getCPUs() (int, error) {
	lscpu, err := cmd.Exec("lscpu", "--json")
	if err != nil {
		fmt.Println(err)
		return 4, nil
		// return 0, err
	}
	var data map[string][]map[string]string
	err = json.Unmarshal(lscpu, &data)
	if err != nil {
		return 0, err
	}

	for _, line := range data["lscpu"] {
		if line["field"] == "CPU(s):" {
			val, err := strconv.Atoi(line["data"])
			return val, err
		}
	}

	return 0, fmt.Errorf("Could not get CPUs")
}

func getLoad() ([]float64, error) {
	uptime, err := cmd.Exec("uptime")
	if err != nil {
		return nil, err
	}
	load := strings.Split(string(uptime), "load average: ")
	loads := strings.Split(load[len(load)-1], ", ")
	vals := make([]float64, len(loads))
	for i, l := range loads {
		parsed, err := strconv.ParseFloat(strings.TrimSpace(l), 64)
		if err != nil {
			return nil, err
		}
		vals[i] = parsed * 100
	}
	return vals, nil
}
