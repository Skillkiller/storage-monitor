package collectors

import (
	"fmt"
	"ssh-storage-monitor/config"
	"ssh-storage-monitor/metrics"
	"ssh-storage-monitor/ssh"

	"github.com/prometheus/client_golang/prometheus"
)

type SshCollector struct {
	hostConfig config.SshHostConfig
}

func NewSshCollector(hostConfig config.SshHostConfig) *SshCollector {
	return &SshCollector{
		hostConfig: hostConfig,
	}
}

func (c *SshCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *SshCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Println("ssh collect")
	connection := ssh.GetSshConnection(c.hostConfig)
	defer connection.Close()

	errorCount := 0

	for _, dirConfig := range c.hostConfig.WatchedDirs {
		size, err := metrics.GetDirSize(connection, dirConfig)
		if err == nil {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					"folder_size",
					"Folder size in bytes",
					[]string{"host", "name"},
					nil,
				),
				prometheus.GaugeValue,
				float64(size),
				c.hostConfig.Host,
				dirConfig.Name)
		} else {
			errorCount++
		}
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"error_status",
			"Indicates if there was an error while collecting metrics",
			[]string{"host"},
			nil,
		),
		prometheus.GaugeValue,
		float64(errorCount),
		c.hostConfig.Host)
}
