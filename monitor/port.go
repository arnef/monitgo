package monitor

type DataReceiver interface {
	Push(data Data)
}

type Status struct {
	Name      string
	Error     *string
	Container map[string]ContainerStats
	Host      HostStats
}

func NewStatusError(name string, err string) Status {
	return Status{Error: &err, Name: name}
}

type HostStats struct {
	CPU    float64
	Memory UsageStats
	Disk   UsageStats
}

type ContainerStats struct {
	Name    string
	CPU     float64
	Memory  UsageStats
	Network NetworkStats
}

type NetworkStats struct {
	RxBytesPerSecond uint64
	TxBytesPerSecond uint64
}

type UsageStats struct {
	TotalBytes uint64
	UsedBytes  uint64
	Percentage float64
}

type NodeConfig struct {
	Name string
	Host string
	Port uint
}

type Data map[string]Status
