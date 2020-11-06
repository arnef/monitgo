package host

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
)

func getNormalizedLoad() ([]int, error) {
	cpus, err := getCPUs()
	if err != nil {
		return nil, err
	}
	load, err := getLoad()
	if err != nil {
		return nil, err
	}
	for i, l := range load {
		load[i] = l / cpus
		// fmt.Println(l / cpus)
	}
	return load, nil
}

func getCPUs() (int, error) {
	lscpu, err := cmd.Exec("lscpu", "--json")
	if err != nil {
		return 0, err
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

func getLoad() ([]int, error) {
	uptime, err := cmd.Exec("uptime")
	if err != nil {
		return nil, err
	}
	load := strings.Split(string(uptime), "load average: ")
	loads := strings.Split(load[len(load)-1], ", ")
	vals := make([]int, len(loads))
	for i, l := range loads {
		parsed, err := strconv.ParseFloat(strings.TrimSpace(l), 64)
		if err != nil {
			return nil, err
		}
		vals[i] = int(parsed * 100)
	}
	return vals, nil
}
