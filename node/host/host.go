package host

type HostStats struct {
	CPULoad   []int
	MemUsage  Usage
	DiskUsage Usage
}

type Usage struct {
	Total      int
	Used       int
	Percentage int
}

func GetStats() (*HostStats, error) {

	cpuLoad, err := getNormalizedLoad()
	if err != nil {
		return nil, err
	}

	memUsage, err := getMemUsage()
	if err != nil {
		return nil, err
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		return nil, err
	}

	var memUsageCom Usage
	for _, mem := range memUsage {
		if mem.Name == "Mem" {
			memUsageCom = Usage{
				Used:       mem.Used,
				Total:      mem.Total,
				Percentage: mem.Used * 100 / mem.Total,
			}
		}
	}

	diskUsageCom := Usage{Used: 0, Total: 0}
	for _, disk := range diskUsage {
		diskUsageCom.Total += disk.Total
		diskUsageCom.Used += disk.Used
	}
	diskUsageCom.Percentage = diskUsageCom.Used * 100 / diskUsageCom.Total

	stats := HostStats{
		CPULoad:   cpuLoad,
		MemUsage:  memUsageCom,
		DiskUsage: diskUsageCom,
	}

	return &stats, nil
}
