package host

type Stats struct {
	CPULoad   []float64
	MemUsage  Usage
	DiskUsage Usage
}

type Usage struct {
	// Total avaiable space in bytes
	Total uint64
	// Used space in bytes
	Used uint64
	// Percentage used * 100 / total
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
