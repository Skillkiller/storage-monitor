package collectors

import (
	"ssh-storage-monitor/config"

	"github.com/prometheus/client_golang/prometheus"
)

type RcloneCollector struct {
	config config.Config
}

func (c *RcloneCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *RcloneCollector) Collect(ch chan<- prometheus.Metric) {

}

func NewRcloneCollector(config config.Config) *RcloneCollector {
	return &RcloneCollector{
		config: config,
	}
}
