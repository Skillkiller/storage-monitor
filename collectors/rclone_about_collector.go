package collectors

import (
	"encoding/json"
	"os/exec"
	"storage-monitor/config"

	"github.com/prometheus/client_golang/prometheus"
)

// CommandExecutor executes an external command and returns its stdout.
type CommandExecutor func(name string, arg ...string) ([]byte, error)

var defaultCommandExecutor CommandExecutor = func(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

// AboutInfo represents the JSON output of `rclone about <remote> --json`.
// Fields are pointers because different rclone backends report different subsets.
type AboutInfo struct {
	Total   *int64 `json:"total"`
	Used    *int64 `json:"used"`
	Trashed *int64 `json:"trashed"`
	Other   *int64 `json:"other"`
	Free    *int64 `json:"free"`
}

// RcloneAboutCollector collects remote storage quota metrics using `rclone about --json`.
type RcloneAboutCollector struct {
	target               config.AboutTarget
	rcloneExecutablePath string
	cmdExecutor          CommandExecutor
	totalBytesDesc       *prometheus.Desc
	usedBytesDesc        *prometheus.Desc
	freeBytesDesc        *prometheus.Desc
	trashedBytesDesc     *prometheus.Desc
	otherBytesDesc       *prometheus.Desc
	upDesc               *prometheus.Desc
}

// NewRcloneAboutCollector creates a collector for an rclone about target.
func NewRcloneAboutCollector(target config.AboutTarget, rcloneExecutablePath string) *RcloneAboutCollector {
	return NewRcloneAboutCollectorWithExecutor(target, rcloneExecutablePath, defaultCommandExecutor)
}

// NewRcloneAboutCollectorWithExecutor creates a collector with a custom CommandExecutor (useful for testing).
func NewRcloneAboutCollectorWithExecutor(target config.AboutTarget, rcloneExecutablePath string, cmdExecutor CommandExecutor) *RcloneAboutCollector {
	labels := []string{"name"}
	return &RcloneAboutCollector{
		target:               target,
		rcloneExecutablePath: rcloneExecutablePath,
		cmdExecutor:          cmdExecutor,
		totalBytesDesc:       prometheus.NewDesc("rclone_about_total_bytes", "Total storage capacity of the remote in bytes", labels, nil),
		usedBytesDesc:        prometheus.NewDesc("rclone_about_used_bytes", "Used storage space on the remote in bytes", labels, nil),
		freeBytesDesc:        prometheus.NewDesc("rclone_about_free_bytes", "Free storage space on the remote in bytes", labels, nil),
		trashedBytesDesc:     prometheus.NewDesc("rclone_about_trashed_bytes", "Trashed storage space on the remote in bytes", labels, nil),
		otherBytesDesc:       prometheus.NewDesc("rclone_about_other_bytes", "Other storage space used on the remote in bytes", labels, nil),
		upDesc:               prometheus.NewDesc("rclone_about_up", "Whether rclone about collection was successful (1 = OK, 0 = Error)", labels, nil),
	}
}

func (c *RcloneAboutCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalBytesDesc
	ch <- c.usedBytesDesc
	ch <- c.freeBytesDesc
	ch <- c.trashedBytesDesc
	ch <- c.otherBytesDesc
	ch <- c.upDesc
}

func (c *RcloneAboutCollector) Collect(ch chan<- prometheus.Metric) {
	out, err := c.cmdExecutor(c.rcloneExecutablePath, "about", c.target.Remote, "--json")
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, c.target.Name)
		return
	}

	var info AboutInfo
	if err := json.Unmarshal(out, &info); err != nil {
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, c.target.Name)
		return
	}

	metrics := []struct {
		desc  *prometheus.Desc
		value *int64
	}{
		{c.totalBytesDesc, info.Total},
		{c.usedBytesDesc, info.Used},
		{c.freeBytesDesc, info.Free},
		{c.trashedBytesDesc, info.Trashed},
		{c.otherBytesDesc, info.Other},
	}

	for _, m := range metrics {
		if m.value != nil {
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(*m.value), c.target.Name)
		}
	}

	ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1, c.target.Name)
}
