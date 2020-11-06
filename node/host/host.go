package host

type Stats struct {
	CPULoad   []float64
	MemUsage  Usage
	DiskUsage Usage
}

type Usage struct {
	Total      float64
	Used       float64
	Percentage float64
}

func GetStats() (*Stats, error) {

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

	stats := Stats{
		CPULoad:   cpuLoad,
		MemUsage:  *memUsage,
		DiskUsage: *diskUsage,
	}

	return &stats, nil
}
