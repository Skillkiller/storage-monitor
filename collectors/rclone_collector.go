package collectors

import (
	"ssh-storage-monitor/config"

	"github.com/prometheus/client_golang/prometheus"
)

type RcloneCollector struct {
	recursiveWatchedDirs config.RecursiveWatchedDirs
}

func (c *RcloneCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *RcloneCollector) Collect(ch chan<- prometheus.Metric) {
}

func NewRcloneCollector(recursiveWatchedDirs config.RecursiveWatchedDirs) *RcloneCollector {
	return &RcloneCollector{
		recursiveWatchedDirs: recursiveWatchedDirs,
	}
}
