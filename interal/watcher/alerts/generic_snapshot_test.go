package alerts_test

import (
	"testing"

	"github.com/arnef/monitgo/interal/watcher/alerts"
	"github.com/arnef/monitgo/pkg"
)

func TestContainerWentDown(t *testing.T) {

	current := alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 0,
			UsedBytes:  0,
		},
		State: pkg.ContainerStateDead,
	}

	if !current.WentDown(nil) {
		t.Error("Expected current to be down")
	}

	if current.WentDown(&alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 0,
			UsedBytes:  0,
		},
		State: pkg.ContainerStateDead,
	}) {
		t.Error("Expected not a new alert because container already was down")
	}

	if !current.WentDown(&alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 1,
			UsedBytes:  0,
		},
		State: pkg.ContainerStateRunning,
	}) {
		t.Error("Expected current to be down")
	}

	if !current.WentDown(&alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 1,
			UsedBytes:  1,
		},
		State: pkg.ContainerStateRunning,
	}) {
		t.Error("Expected current to be down")
	}
	// not down just idyling
	current = alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 1,
			UsedBytes:  0,
		},
		State: pkg.ContainerStateRunning,
	}

	if current.WentDown(nil) {
		t.Error("Expect current not to be down")
	}
}

func TestContainerCameUp(t *testing.T) {

	current := alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 1,
			UsedBytes:  1,
		},
		State: pkg.ContainerStateRunning,
	}
	if !current.CameUp(&alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 0,
			UsedBytes:  0,
		},
		State: pkg.ContainerStateDead,
	}) {
		t.Error("current came up")
	}

	if current.CameUp(&alerts.GenericSnaphot{
		MemoryUsage: &pkg.Usage{
			TotalBytes: 1,
			UsedBytes:  1,
		},
		State: pkg.ContainerStateRunning,
	}) {
		t.Error("current does not came up")
	}

	if current.CameUp(nil) {
		t.Error("current does not came up (but started)")
	}
}
